package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/lib/pq"
	"github.com/pquerna/otp/totp"
	"image/png"
	"time"
)

var (
	//ErrWrongOTPCode = errors.New("wrong OTP provided")
	ErrUnAuthorized = errors.New("Not authorized")
	ErrBadRequest   = errors.New("Bad request")
	//ErrServer       = errors.New("Internal Server Error")
)

type UserService interface {
	Register(ctx context.Context, token *token.Payload, req dto.CreateUserRequest) (dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginUserRequest) (response dto.LoginUserResponse, roles []int64, schoolID int64, ClassID int64, err error)
	CreateAdmin() error
	CheckTOTP(ctx context.Context, token *token.Payload, req dto.CheckTOTPRequest) (response dto.LoginUserResponse, err error)
}

type userService struct {
	RolesService
	ClassService
	db    db.Store
	roles map[string]db.Role
}

// TODO uuid Roles should be hidden by uuid

func (s *userService) Register(ctx context.Context, authToken *token.Payload, req dto.CreateUserRequest) (dto.UserResponse, error) {
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
	_, err = s.db.GetRolebyId(ctx, req.RoleID)
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

	// Check that the class is provided only if request role is Student
	if ((req.RoleID == s.roles[Student].ID) && req.ClassID == 0) || ((req.RoleID != s.roles[Student].ID) && req.ClassID != 0) {
		return dto.UserResponse{}, ErrBadRequest
	}

	hashedPassword, err := HashPassword(req.Password)
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
	response.TOTPSecret = user.TotpSecret
	response.UserClass = userRoleClass.ClassID
	response.UserRole = userRole.RoleID
	response.UserSchool = school.ID
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

func NewUserService(database db.Store, mapRoles map[string]db.Role) UserService {

	return &userService{db: database, roles: mapRoles,
		RolesService: NewRolesService(database),
		ClassService: NewClassService(database, mapRoles)}
}
