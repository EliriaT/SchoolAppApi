package service

import (
	"context"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
)

type Service interface {
	UserService
	SchoolService
}

type ServerService struct {
	UserService
	SchoolService
}

func NewServerService(database db.Store) (Service, error) {

	userRoles, err := database.GetRoles(context.TODO())
	if err != nil {
		return &ServerService{}, err
	}
	mapRoles := make(map[string]db.Role)
	for _, r := range userRoles {
		mapRoles[r.Name] = r
	}

	return &ServerService{UserService: NewUserService(database, mapRoles),
		SchoolService: NewSchoolService(database)}, err
}
