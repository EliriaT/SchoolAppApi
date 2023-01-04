package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

type ClassService interface {
	AddUserToClass(ctx context.Context, userRoleID int64, classID int64) (db.UserRoleClass, error)
	GetUserClassByUserRoleId(ctx context.Context, userRoleID int64) (db.UserRoleClass, error)
	GetClassByID(ctx context.Context, authToken *token.Payload, classID int64) (dto.ClassResponse, error)
	GetClass(ctx context.Context, authToken *token.Payload) ([]dto.ClassResponse, error)
	CreateClass(ctx context.Context, authToken *token.Payload, req dto.CreateClassRequest) (dto.ClassResponse, error)
	ChangeHeadTeacherClass(ctx context.Context, authToken *token.Payload, req dto.ChangeHeadTeacherRequest) (dto.ClassResponse, error)
}

type classService struct {
	db db.Store
	RolesService
	roles map[string]db.Role
}

// called in user service
func (rs *classService) AddUserToClass(ctx context.Context, userRoleID int64, classID int64) (db.UserRoleClass, error) {
	args := db.AddUserToClassParams{
		UserRoleID: userRoleID,
		ClassID:    classID,
	}
	userRoleClass, err := rs.db.AddUserToClass(ctx, args)
	return userRoleClass, err
}

// called in user service
func (rs *classService) GetUserClassByUserRoleId(ctx context.Context, userRoleID int64) (db.UserRoleClass, error) {
	userRoleClass, err := rs.db.GetUserClassByUserRoleId(ctx, userRoleID)
	return userRoleClass, err
}

// called in handler
func (rs *classService) GetClassByID(ctx context.Context, authToken *token.Payload, classID int64) (dto.ClassResponse, error) {
	if !CheckRolePresence(authToken.Role, rs.roles[Director].ID) && !CheckRolePresence(authToken.Role, rs.roles[SchoolManager].ID) && !CheckRolePresence(authToken.Role, rs.roles[HeadTeacher].ID) {
		return dto.ClassResponse{}, ErrUnAuthorized
	}
	if CheckRolePresence(authToken.Role, rs.roles[HeadTeacher].ID) {
		classID = authToken.ClassID
	}

	class, err := rs.db.GetClassById(ctx, classID)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	students, err := rs.db.GetClassWithStudents(ctx, classID)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	var headTeacherName string
	pupils := make([]dto.UserResponse, 0, 35)
	for _, s := range students {
		if s.RoleID == rs.roles[HeadTeacher].ID {
			headTeacherName = s.FirstName + s.LastName
			continue
		}
		student := dto.UserResponse{
			Email:       s.Email,
			LastName:    s.LastName,
			FirstName:   s.FirstName,
			Gender:      s.Gender,
			PhoneNumber: s.PhoneNumber,
			Domicile:    s.Domicile,
			BirthDate:   s.BirthDate,
		}
		pupils = append(pupils, student)
	}

	response := dto.ClassResponse{
		ID:              class.ID,
		Name:            class.Name,
		HeadTeacherName: headTeacherName,
		Pupils:          pupils,
	}

	return response, err
}

// the id is the userRole ID of the teacher /director / school manager. I should add another user role as Head Teacher, and Add in UserRoleClass the info
func (rs *classService) CreateClass(ctx context.Context, authToken *token.Payload, req dto.CreateClassRequest) (dto.ClassResponse, error) {
	if !CheckRolePresence(authToken.Role, rs.roles[Director].ID) && !CheckRolePresence(authToken.Role, rs.roles[SchoolManager].ID) {
		return dto.ClassResponse{}, ErrUnAuthorized
	}

	userRole, err := rs.db.GetUserRoleById(ctx, req.HeadTeacher)
	if err != nil {
		return dto.ClassResponse{}, err
	}
	createUserRoleArg := db.CreateRoleForUserParams{RoleID: rs.roles[HeadTeacher].ID, UserID: userRole.UserID, SchoolID: authToken.SchoolID}
	createdUserRole, err := rs.db.CreateRoleForUser(ctx, createUserRoleArg)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	// here the role is from teacher, but it anyway maps to the same user
	arg := db.CreateClassParams{Name: req.Name, HeadTeacher: createdUserRole.ID}
	class, err := rs.db.CreateClass(ctx, arg)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	addUserClassArgs := db.AddUserToClassParams{ClassID: class.ID, UserRoleID: createdUserRole.ID}
	_, err = rs.db.AddUserToClass(ctx, addUserClassArgs)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	response := dto.ClassResponse{ID: class.ID,
		Name:        class.Name,
		HeadTeacher: class.HeadTeacher}
	return response, err
}

// if role is school manager or director, then a list of classes as response, otherwise his/her class
func (rs *classService) GetClass(ctx context.Context, authToken *token.Payload) ([]dto.ClassResponse, error) {
	if CheckRolePresence(authToken.Role, rs.roles[Admin].ID) {
		return []dto.ClassResponse{}, ErrUnAuthorized
	}
	if CheckRolePresence(authToken.Role, rs.roles[Director].ID) || CheckRolePresence(authToken.Role, rs.roles[SchoolManager].ID) {
		classes, err := rs.db.ListAllClasses(ctx)
		if err != nil {
			return nil, err
		}
		response := make([]dto.ClassResponse, 0, 50)

		for _, class := range classes {
			teacherUser, err := rs.GetUserFromUserRoleID(ctx, class.HeadTeacher)
			if err != nil {
				return nil, err
			}
			teacherName := teacherUser.LastName + " " + teacherUser.FirstName
			newClass := dto.ClassResponse{
				ID:              class.ID,
				Name:            class.Name,
				HeadTeacher:     class.HeadTeacher,
				HeadTeacherName: teacherName,
			}
			response = append(response, newClass)
		}
		return response, err
	}
	if CheckRolePresence(authToken.Role, rs.roles[HeadTeacher].ID) || CheckRolePresence(authToken.Role, rs.roles[Student].ID) {
		classId := authToken.ClassID
		class, err := rs.db.GetClassById(ctx, classId)
		if err != nil {
			return []dto.ClassResponse{}, err
		}

		students, err := rs.db.GetClassWithStudents(ctx, classId)
		if err != nil {
			return []dto.ClassResponse{}, err
		}

		var headTeacherName string
		pupils := make([]dto.UserResponse, 0, 35)
		for _, s := range students {
			if s.RoleID == rs.roles[HeadTeacher].ID {
				headTeacherName = s.FirstName + " " + s.LastName
				continue
			}
			var student dto.UserResponse
			if CheckRolePresence(authToken.Role, rs.roles[HeadTeacher].ID) {
				student = dto.UserResponse{
					Email:       s.Email,
					LastName:    s.LastName,
					FirstName:   s.FirstName,
					Gender:      s.Gender,
					PhoneNumber: s.PhoneNumber,
					Domicile:    s.Domicile,
					BirthDate:   s.BirthDate,
				}
			} else {
				student = dto.UserResponse{
					LastName:  s.LastName,
					FirstName: s.FirstName,
					BirthDate: s.BirthDate,
				}
			}

			pupils = append(pupils, student)
		}

		response := []dto.ClassResponse{{
			ID:              class.ID,
			Name:            class.Name,
			HeadTeacherName: headTeacherName,
			Pupils:          pupils,
		},
		}

		return response, err
	}
	//Happens if it is a teacher
	return []dto.ClassResponse{}, ErrUnAuthorized
}

// TODO update old user_role_class, delete user role
// TODO CHECK THAT REQ.HEADTEACHER IS NOT STUDENT
func (rs *classService) ChangeHeadTeacherClass(ctx context.Context, authToken *token.Payload, req dto.ChangeHeadTeacherRequest) (dto.ClassResponse, error) {
	if !CheckRolePresence(authToken.Role, rs.roles[Director].ID) && !CheckRolePresence(authToken.Role, rs.roles[SchoolManager].ID) {
		return dto.ClassResponse{}, ErrUnAuthorized
	}

	newTeacherUserRole, err := rs.db.GetUserRoleById(ctx, req.HeadTeacherID)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	class, err := rs.db.GetClassById(ctx, req.ClassID)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	oldHeadTeacherUserRole, err := rs.db.GetUserRoleById(ctx, class.HeadTeacher)
	_, err = rs.db.UpdateUserHeadTeacherRole(ctx, db.UpdateUserHeadTeacherRoleParams{UserID: newTeacherUserRole.UserID, ID: oldHeadTeacherUserRole.ID})
	if err != nil {
		return dto.ClassResponse{}, err
	}

	response := dto.ClassResponse{
		ID:          class.ID,
		Name:        class.Name,
		HeadTeacher: class.HeadTeacher,
	}
	return response, err
}
func NewClassService(database db.Store, mapRoles map[string]db.Role) ClassService {
	return &classService{db: database, roles: mapRoles, RolesService: NewRolesService(database, mapRoles)}
}
