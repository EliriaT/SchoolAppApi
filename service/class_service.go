package service

import (
	"context"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
)

type ClassService interface {
	AddUserToClass(ctx context.Context, userRoleID int64, classID int64) (db.UserRoleClass, error)
	GetUserClassByUserRoleId(ctx context.Context, userRoleID int64) (db.UserRoleClass, error)
	GetClass(ctx context.Context, classID int64) (db.Class, error)
}

type classService struct {
	db    db.Store
	roles map[string]db.Role
}

func (rs *classService) AddUserToClass(ctx context.Context, userRoleID int64, classID int64) (db.UserRoleClass, error) {
	args := db.AddUserToClassParams{
		UserRoleID: userRoleID,
		ClassID:    classID,
	}
	userRoleClass, err := rs.db.AddUserToClass(ctx, args)
	return userRoleClass, err
}

func (rs *classService) GetUserClassByUserRoleId(ctx context.Context, userRoleID int64) (db.UserRoleClass, error) {
	userRoleClass, err := rs.db.GetUserClassByUserRoleId(ctx, userRoleID)
	return userRoleClass, err
}

func (rs *classService) GetClass(ctx context.Context, classID int64) (db.Class, error) {
	class, err := rs.db.GetClassById(ctx, classID)
	if err != nil {
		return db.Class{}, err
	}
	return class, err
}
func NewClassService(database db.Store, mapRoles map[string]db.Role) ClassService {
	return &classService{db: database, roles: mapRoles}
}
