package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
)

type Service interface {
	UserService
	SchoolService
	ClassService
	RolesService
	SemesterService
	CourseService
	LessonService
	MarkService
	EmailService
	SessionService
}

type ServerService struct {
	UserService
	SchoolService
	ClassService
	RolesService
	SemesterService
	CourseService
	LessonService
	MarkService
	EmailService
	SessionService
}

func NewServerService(database db.Store, configSet config.Config) (Service, error) {

	userRoles, err := database.GetRoles(context.TODO())
	if err != nil {
		return &ServerService{}, err
	}
	mapRoles := make(map[string]db.Role)
	for _, r := range userRoles {
		mapRoles[r.Name] = r
	}

	return &ServerService{UserService: NewUserService(database, mapRoles, configSet),
		SchoolService:   NewSchoolService(database, mapRoles),
		RolesService:    NewRolesService(database, mapRoles),
		ClassService:    NewClassService(database, mapRoles),
		SemesterService: NewSemesterService(database, mapRoles),
		CourseService:   NewCourseService(database, mapRoles),
		LessonService:   NewLessonService(database, mapRoles),
		MarkService:     NewMarkService(database, mapRoles),
		EmailService:    NewEmailService(configSet.EmailServerLogin, configSet.EmailServerPassword),
		SessionService:  NewSessionService(database)}, err
}
