package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

const (
	Admin         = "Admin"
	Director      = "Director"
	SchoolManager = "School_Manager"
	Teacher       = "Teacher"
	HeadTeacher   = "Head_Teacher"
	Student       = "Student"
)

func CheckRolePresence(roles []int64, checkRole int64) bool {
	for _, r := range roles {
		if r == checkRole {
			return true
		}
	}
	return false
}

type RolesService interface {
	//GetRoles(ctx context.Context, req dto.CreateSchoolRequest) (dto.SchoolResponse, error)
	AddUserRole(ctx context.Context, userID int64, roleID int64, schoolID int64) (db.UserRole, error)
	GetUserFromUserRoleID(ctx context.Context, userRoleID int64) (db.User, error)
	GetRoles(ctx context.Context, authToken *token.Payload) (response []dto.RoleResponse, err error)
}

type roleService struct {
	db    db.Store
	roles map[string]db.Role
}

func (rs *roleService) AddUserRole(ctx context.Context, userID int64, roleID int64, schoolID int64) (db.UserRole, error) {
	args := db.CreateRoleForUserParams{RoleID: roleID,
		UserID:   userID,
		SchoolID: schoolID,
	}
	userRole, err := rs.db.CreateRoleForUser(ctx, args)
	return userRole, err
}

func (rs *roleService) GetUserFromUserRoleID(ctx context.Context, userRoleID int64) (db.User, error) {
	userRole, err := rs.db.GetUserRoleById(ctx, userRoleID)
	if err != nil {
		return db.User{}, err
	}
	user, err := rs.db.GetUserbyId(ctx, userRole.UserID)
	return user, err
}
func (rs *roleService) GetRoles(ctx context.Context, authToken *token.Payload) (response []dto.RoleResponse, err error) {
	if !CheckRolePresence(authToken.Role, rs.roles[Director].ID) && !CheckRolePresence(authToken.Role, rs.roles[HeadTeacher].ID) && !CheckRolePresence(authToken.Role, rs.roles[SchoolManager].ID) {
		return nil, ErrUnAuthorized
	}
	roles, err := rs.db.GetRoles(ctx)
	for _, r := range roles {
		response = append(response, dto.RoleResponse{
			r.ID, r.Name,
		})
	}
	return
}

func NewRolesService(database db.Store, roles map[string]db.Role) RolesService {
	return &roleService{db: database, roles: roles}
}
