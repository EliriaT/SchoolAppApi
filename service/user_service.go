package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/sethvargo/go-password/password"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"io"
	"unicode"

	"github.com/lib/pq"
	"github.com/pquerna/otp/totp"
	"image/png"
	"log"
	"net/http"
	"time"
)

var (
	//ErrWrongOTPCode = errors.New("wrong OTP provided")
	ErrUnAuthorized = errors.New("Not authorized")
	ErrBadRequest   = errors.New("Bad request")
	ErrInvalidEmail = errors.New("Invalid email provided")
	ErrEasyPassword = errors.New("The password is not enough complex")
)

const minEntropy = 60

type UserService interface {
	Register(ctx context.Context, token *token.Payload, req dto.CreateUserRequest, tokenPassReset string) (dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginUserRequest) (response dto.LoginUserResponse, roles []int64, schoolID int64, ClassID int64, err error)
	CreateAdmin() error
	CheckTOTP(ctx context.Context, token *token.Payload, req dto.CheckTOTPRequest) (response dto.LoginUserResponse, err error)
	GetTeachers(ctx context.Context, token *token.Payload) (response []dto.TeacherResponse, err error)
	// this will change the password
	ChangePassword(ctx context.Context, email string, password string) (err error)

	VerifyEmail(email string) (err error)
	ValidatePasswords(password string) error
	CreatePasswordRecoveryLink(token string, email string) (link string)

	CreateSession(ctx context.Context, refreshToken string, refreshPayload *token.Payload, ipAddress string, userAgent string) (dto.SessionResponse, error)
}

type userService struct {
	RolesService
	ClassService
	EmailService
	db            db.Store
	roles         map[string]db.Role
	emailVerifier *emailverifier.Verifier
	configSet     config.Config
}

// TODO uuid Roles should be hidden by uuid

func (s *userService) Register(ctx context.Context, authToken *token.Payload, req dto.CreateUserRequest, tokenPassReset string) (dto.UserResponse, error) {
	// check to see that the role is Admin or Director or School_Manager
	if !CheckRolePresence(authToken.Role, s.roles[Admin].ID) && !CheckRolePresence(authToken.Role, s.roles[Director].ID) && !CheckRolePresence(authToken.Role, s.roles[SchoolManager].ID) && !CheckRolePresence(authToken.Role, s.roles[HeadTeacher].ID) {
		return dto.UserResponse{}, ErrUnAuthorized
	}

	// check to see that the school is provided only if the user is Admin
	if (CheckRolePresence(authToken.Role, s.roles[Admin].ID) && req.SchoolID == 0) || (!(CheckRolePresence(authToken.Role, s.roles[Admin].ID)) && req.SchoolID != 0) {
		return dto.UserResponse{}, ErrBadRequest
	}

	var schoolID int64

	if req.SchoolID == 0 {
		schoolID = authToken.SchoolID
	} else {
		schoolID = req.SchoolID
	}

	// Check that the school is present in the database
	school, err := s.db.GetSchoolbyId(ctx, schoolID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Check that the class is present in the database
	if req.RoleID == s.roles[Student].ID {
		_, err := s.db.GetClassById(ctx, req.ClassID)
		if err != nil {
			return dto.UserResponse{}, err
		}
	}

	// Check the assigned role
	roleEntity, err := s.db.GetRolebyId(ctx, req.RoleID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Admin -> Director or School Manager
	if CheckRolePresence(authToken.Role, s.roles[Admin].ID) && (req.RoleID != s.roles[Director].ID && req.RoleID != s.roles[SchoolManager].ID) {
		return dto.UserResponse{}, ErrUnAuthorized
	}
	// Director and School Manager -> Teachers, Students; Director -> SchoolManager
	if (CheckRolePresence(authToken.Role, s.roles[Director].ID) || CheckRolePresence(authToken.Role, s.roles[SchoolManager].ID)) && (req.RoleID != s.roles[Teacher].ID && req.RoleID != s.roles[Student].ID && req.RoleID != s.roles[SchoolManager].ID) {
		return dto.UserResponse{}, ErrUnAuthorized
	}
	// Director -> School Manager
	if CheckRolePresence(authToken.Role, s.roles[SchoolManager].ID) && req.RoleID == s.roles[SchoolManager].ID {
		return dto.UserResponse{}, ErrUnAuthorized
	}

	//Head Teacher -> student in his class
	if (CheckRolePresence(authToken.Role, s.roles[HeadTeacher].ID)) && (req.RoleID != s.roles[Student].ID) {
		return dto.UserResponse{}, ErrUnAuthorized
	}

	//Head Teacher -> student in his class
	if (CheckRolePresence(authToken.Role, s.roles[HeadTeacher].ID)) && (req.RoleID == s.roles[Student].ID) {
		if authToken.ClassID != req.ClassID {
			return dto.UserResponse{}, ErrUnAuthorized
		}

	}

	// Check that the class is provided only if request role is Student
	if ((req.RoleID == s.roles[Student].ID) && req.ClassID == 0) || ((req.RoleID != s.roles[Student].ID) && req.ClassID != 0) {
		return dto.UserResponse{}, ErrBadRequest
	}

	securePassword, err := password.Generate(12, 5, 3, false, true)
	if err != nil {
		return dto.UserResponse{}, err
	}
	hashedPassword, err := HashPassword(securePassword)
	if err != nil {
		return dto.UserResponse{}, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ESchool.md",
		AccountName: req.Email,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	arg := db.CreateUserParams{
		Email:      req.Email,
		Password:   hashedPassword,
		TotpSecret: key.Secret(),
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		Gender:     req.Gender,
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  true,
		},
		Domicile: sql.NullString{
			String: req.Domicile,
			Valid:  true,
		},
		BirthDate: sql.NullTime{
			Time:  req.BirthDate,
			Valid: true,
		},
	}

	user, err := s.db.CreateUser(ctx, arg)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Add Role to User in db
	userRole, err := s.AddUserRole(ctx, user.ID, req.RoleID, school.ID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	var userRoleClass db.UserRoleClass
	// Add user to class (if it is head teacher or student
	if req.RoleID == s.roles[Student].ID || req.RoleID == s.roles[HeadTeacher].ID {
		userRoleClass, err = s.AddUserToClass(ctx, userRole.ID, req.ClassID)
		if err != nil {
			return dto.UserResponse{}, err
		}
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	err = png.Encode(&buf, img)
	if err != nil {
		return dto.UserResponse{}, err
	}
	qrimage := base64.StdEncoding.EncodeToString(buf.Bytes())
	response := dto.NewUserResponse(user)
	response.Qrcode = qrimage
	//response.TOTPSecret = user.TotpSecret
	response.UserClass = userRoleClass.ClassID
	response.UserRole = userRole.RoleID
	response.UserSchool = school.ID

	link := s.CreatePasswordRecoveryLink(tokenPassReset, user.Email)
	err = s.SendRegisterEmail(user.FirstName+" "+user.LastName, school.Name, roleEntity.Name, securePassword, qrimage, user.Email, link)
	if err != nil {
		return dto.UserResponse{}, ErrSendingEmail
	}
	return response, err
}

func (s *userService) Login(ctx context.Context, req dto.LoginUserRequest) (response dto.LoginUserResponse, roles []int64, schoolID int64, ClassID int64, err error) {
	user, err := s.db.GetUserbyEmail(ctx, req.Email)
	if err != nil {
		return
	}

	err = CheckPassword(req.Password, user.Password)
	if err != nil {
		return
	}

	userRoles, err := s.db.GetUserRoleByUserId(ctx, user.ID)
	if err != nil {
		return
	}

	for _, ur := range userRoles {
		roles = append(roles, ur.RoleID)
	}
	schoolID = userRoles[0].SchoolID

	for i, value := range roles {
		if value == s.roles[HeadTeacher].ID || value == s.roles[Student].ID {
			userClass, err := s.GetUserClassByUserRoleId(ctx, userRoles[i].ID)
			if err != nil {
				return response, roles, schoolID, ClassID, err
			}
			ClassID = userClass.ClassID
			break
		}
	}

	response = dto.LoginUserResponse{User: dto.NewUserResponse(user)}
	return

}

// here should be environmental variables
func (s *userService) CreateAdmin() error {
	school, err := s.db.CreateSchool(context.TODO(), "admin")
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil
			default:
				return err
			}
		}
		return err
	}

	hashedPassword, err := HashPassword("123456")
	if err != nil {
		return err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ESchool.md",
		AccountName: "irinaAdmin@gmail.com",
		Secret:      []byte("JNHMDMDPX6RY3Z3CJBB6AHTY6BDOQDDY"),
	})
	if err != nil {
		return err
	}

	arg := db.CreateUserParams{
		Email:      "irinaAdmin@gmail.com",
		Password:   hashedPassword,
		TotpSecret: key.Secret(),
		LastName:   "Tiora",
		FirstName:  "Irina",
		Gender:     "F",
		PhoneNumber: sql.NullString{
			String: "078111111",
			Valid:  true,
		},
		Domicile: sql.NullString{
			String: "'Str. Meaw mur ",
			Valid:  true,
		},
		BirthDate: sql.NullTime{
			Time:  time.Date(2011, time.November, 11, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	}

	user, err := s.db.CreateUser(context.TODO(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil
			default:
				return err
			}
		}
		return err
	}

	_, err = s.AddUserRole(context.TODO(), user.ID, s.roles[Admin].ID, school.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil
			default:
				return err
			}
		}
		return err
	}
	return err
}

func (s *userService) CheckTOTP(ctx context.Context, token *token.Payload, req dto.CheckTOTPRequest) (response dto.LoginUserResponse, err error) {
	user, err := s.db.GetUserbyId(ctx, token.UserID)
	if err != nil {
		return
	}

	valid := totp.Validate(req.TotpToken, user.TotpSecret)
	if !valid {
		err = ErrUnAuthorized
		return
	}
	response = dto.LoginUserResponse{User: dto.NewUserResponse(user)}
	return
}

func (s *userService) GetTeachers(ctx context.Context, token *token.Payload) (response []dto.TeacherResponse, err error) {
	if !CheckRolePresence(token.Role, s.roles[Director].ID) && !CheckRolePresence(token.Role, s.roles[SchoolManager].ID) {
		return []dto.TeacherResponse{}, ErrUnAuthorized
	}
	teachers, err := s.db.GetTeachers(ctx, token.SchoolID)
	if err != nil {
		return []dto.TeacherResponse{}, err
	}
	log.Println("ajunge")
	for _, t := range teachers {
		response = append(response, dto.TeacherResponse{UserID: t.UserID, UserName: t.FirstName + " " + t.LastName, UserRoleID: t.ID_2, RoleName: t.Name})
	}
	return
}

func (s *userService) ChangePassword(ctx context.Context, email string, password string) (err error) {
	user, err := s.db.GetUserbyEmail(ctx, email)
	if err != nil {
		return err
	}
	err = s.ValidatePasswords(password)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		Password: hashedPassword,
		ID:       user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CreatePasswordRecoveryLink(token string, email string) (link string) {
	return s.configSet.ServerAddress + "/users/accountrecovery/" + email + "/" + token
}

// not a disposable email, and there is a user with such an email
func (s *userService) VerifyEmail(email string) (err error) {
	result, err := s.emailVerifier.Verify(email)
	if err != nil {
		return err
	}
	if !result.Syntax.Valid || result.Disposable || result.Suggestion != "" || result.Reachable == "no" || !result.HasMxRecords {
		return ErrInvalidEmail
	}
	_, err = s.db.GetUserbyEmail(context.TODO(), email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func (s *userService) ValidatePasswords(password string) error {
	err := passwordvalidator.Validate(password, minEntropy)
	if err != nil {
		return ErrEasyPassword
	}
	var flag = 0
	for _, char := range password {
		if !unicode.IsLetter(char) {
			flag++
		}
	}
	if flag < 2 {
		return ErrEasyPassword
	}
	return nil
}

func (s *userService) CreateSession(ctx context.Context, refreshToken string, refreshPayload *token.Payload, ipAddress string, userAgent string) (dto.SessionResponse, error) {
	session, err := s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        refreshPayload.Email,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     ipAddress,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return dto.NewSessionResponse(session), err
	}
	return dto.NewSessionResponse(session), nil
}

func NewUserService(database db.Store, mapRoles map[string]db.Role, configSet config.Config) UserService {
	ev := setupEmailVerifier()

	return &userService{db: database, roles: mapRoles,
		RolesService:  NewRolesService(database, mapRoles),
		ClassService:  NewClassService(database, mapRoles),
		EmailService:  NewEmailService(configSet.EmailServerLogin, configSet.EmailServerPassword),
		emailVerifier: ev,
		configSet:     configSet}
}

func setupEmailVerifier() *emailverifier.Verifier {
	ev := emailverifier.NewVerifier()
	ev.EnableDomainSuggest()
	resp, err := http.Get("https://disposable.github.io/disposable-email-domains/domains.json")
	if err != nil {
		log.Println("Error Downloading disposable email domains list")
	}
	defer resp.Body.Close()
	var disposableEmailDomains = make([]string, 0, 5000)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(respBody, &disposableEmailDomains)
	if err != nil {
		log.Println(err)
	}

	ev.AddDisposableDomains(disposableEmailDomains)
	return ev
}
