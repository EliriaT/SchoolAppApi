package service

import (
	"context"
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

type MarkService interface {
	CreateMark(ctx context.Context, authToken *token.Payload, req dto.CreateMarkRequest) (dto.MarkResponse, error)

	ChangeMark(ctx context.Context, authToken *token.Payload, req dto.UpdateMarkRequest) (dto.MarkResponse, error)

	DeleteMark(ctx context.Context, authToken *token.Payload, req dto.DeleteMarkRequest) error
}

type markService struct {
	db    db.Store
	roles map[string]db.Role
}

func (m *markService) CreateMark(ctx context.Context, authToken *token.Payload, req dto.CreateMarkRequest) (dto.MarkResponse, error) {
	if !CheckRolePresence(authToken.Role, m.roles[Teacher].ID) && !CheckRolePresence(authToken.Role, m.roles[Director].ID) && !CheckRolePresence(authToken.Role, m.roles[SchoolManager].ID) {
		return dto.MarkResponse{}, ErrUnAuthorized
	}
	// can create marks only the staff that is teaching the subject!!!
	userRoles, err := m.db.GetUserRoleByUserId(ctx, authToken.UserID)
	if err != nil {
		return dto.MarkResponse{}, err
	}
	course, err := m.db.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return dto.MarkResponse{}, err
	}
	valid := false
	for _, ur := range userRoles {
		if course.TeacherID == ur.ID {
			valid = true
		}
	}
	if !valid {
		return dto.MarkResponse{}, ErrUnAuthorized
	}
	mark, err := m.db.CreateMark(ctx, db.CreateMarkParams{Mark: sql.NullInt32{req.Mark, true}, MarkDate: req.MarkDate, IsAbsent: sql.NullBool{req.IsAbsent, true}, CourseID: req.CourseID, StudentID: req.StudentID})
	if err != nil {
		return dto.MarkResponse{}, err
	}
	return dto.NewMarkResponse(mark), err
}

func (m *markService) ChangeMark(ctx context.Context, authToken *token.Payload, req dto.UpdateMarkRequest) (dto.MarkResponse, error) {
	if !CheckRolePresence(authToken.Role, m.roles[Teacher].ID) && !CheckRolePresence(authToken.Role, m.roles[Director].ID) && !CheckRolePresence(authToken.Role, m.roles[SchoolManager].ID) {
		return dto.MarkResponse{}, ErrUnAuthorized
	}
	// can update marks only the staff that is teaching the subject!!!
	userRoles, err := m.db.GetUserRoleByUserId(ctx, authToken.UserID)
	if err != nil {
		return dto.MarkResponse{}, err
	}
	course, err := m.db.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return dto.MarkResponse{}, err
	}
	valid := false
	for _, ur := range userRoles {
		if course.TeacherID == ur.ID {
			valid = true
		}
	}
	if !valid {
		return dto.MarkResponse{}, ErrUnAuthorized
	}

	if req.Mark == 0 && req.IsAbsent == false || req.Mark != 0 && req.IsAbsent == true {
		return dto.MarkResponse{}, ErrBadRequest
	}
	var mark db.Mark
	if req.Mark != 0 {
		mark, err = m.db.UpdateCourseMarksbyId(ctx, db.UpdateCourseMarksbyIdParams{Mark: sql.NullInt32{req.Mark, true}, ID: req.MarkID})
		if err != nil {
			return dto.MarkResponse{}, err
		}
	}
	if req.IsAbsent == true {
		mark, err = m.db.UpdateCourseAbsencebyId(ctx, req.MarkID)
		if err != nil {
			return dto.MarkResponse{}, err
		}
	}
	return dto.NewMarkResponse(mark), err

}

func (m *markService) DeleteMark(ctx context.Context, authToken *token.Payload, req dto.DeleteMarkRequest) error {
	if !CheckRolePresence(authToken.Role, m.roles[Teacher].ID) && !CheckRolePresence(authToken.Role, m.roles[Director].ID) && !CheckRolePresence(authToken.Role, m.roles[SchoolManager].ID) {
		return ErrUnAuthorized
	}
	// can update marks only the staff that is teaching the subject!!!
	userRoles, err := m.db.GetUserRoleByUserId(ctx, authToken.UserID)
	if err != nil {
		return err
	}

	mark, err := m.db.GetMarkByID(ctx, req.MarkID)
	if err != nil {
		return err
	}

	course, err := m.db.GetCourseByID(ctx, mark.CourseID)
	if err != nil {
		return err
	}
	valid := false
	for _, ur := range userRoles {
		if course.TeacherID == ur.ID {
			valid = true
		}
	}
	if !valid {
		return ErrUnAuthorized
	}
	err = m.db.DeleteMark(ctx, req.MarkID)
	return err
}

func NewMarkService(database db.Store, mapRoles map[string]db.Role) MarkService {
	return &markService{db: database, roles: mapRoles}
}
