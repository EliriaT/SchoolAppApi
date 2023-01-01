package service

import (
	"context"
	"database/sql"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
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
}

type roleService struct {
	db db.Store
}

func (rs *roleService) AddUserRole(ctx context.Context, userID int64, roleID int64, schoolID int64) (db.UserRole, error) {
	args := db.CreateRoleForUserParams{RoleID: roleID,
		UserID: userID,
		SchoolID: sql.NullInt64{
			Int64: schoolID,
			Valid: true,
		}}
	userRole, err := rs.db.CreateRoleForUser(ctx, args)
	return userRole, err
}

func NewRolesService(database db.Store) RolesService {
	return &roleService{db: database}
}
