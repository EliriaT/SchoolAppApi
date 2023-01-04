package service

import (
	"context"
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
		SchoolService:   NewSchoolService(database, mapRoles),
		RolesService:    NewRolesService(database),
		ClassService:    NewClassService(database, mapRoles),
		SemesterService: NewSemesterService(database, mapRoles),
		CourseService:   NewCourseService(database, mapRoles),
		LessonService:   NewLessonService(database, mapRoles),
		MarkService:     NewMarkService(database, mapRoles)}, err
}
