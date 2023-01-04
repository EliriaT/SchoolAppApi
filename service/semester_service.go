package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

type SemesterService interface {
	CreateSemester(ctx context.Context, authToken *token.Payload, req dto.CreateSemesterRequest) (dto.SemesterResponse, error)
	ListSemesters(ctx context.Context, authToken *token.Payload) (response []dto.SemesterResponse, err error)
	GetCurrentSemester(ctx context.Context, authToken *token.Payload) (dto.SemesterResponse, error)
}

type semesterService struct {
	db    db.Store
	roles map[string]db.Role
}

func (s *semesterService) CreateSemester(ctx context.Context, authToken *token.Payload, req dto.CreateSemesterRequest) (dto.SemesterResponse, error) {

	if !CheckRolePresence(authToken.Role, s.roles[Director].ID) && !CheckRolePresence(authToken.Role, s.roles[SchoolManager].ID) {
		return dto.SemesterResponse{}, ErrUnAuthorized
	}

	semester, err := s.db.CreateSemester(ctx, db.CreateSemesterParams{Name: req.Name, StartDate: req.StartDate, EndDate: req.EndDate})
	if err != nil {
		return dto.SemesterResponse{}, err
	}

	response := dto.SemesterResponse{
		ID:        semester.ID,
		Name:      semester.Name,
		StartDate: semester.StartDate,
		EndDate:   semester.EndDate,
		CreatedBy: semester.CreatedBy.Int64,
		UpdatedBy: semester.UpdatedBy.Int64,
		CreatedAt: semester.CreatedAt.Time,
		UpdatedAt: semester.UpdatedAt.Time,
	}
	return response, err
}

func (s *semesterService) ListSemesters(ctx context.Context, authToken *token.Payload) (response []dto.SemesterResponse, err error) {

	if !CheckRolePresence(authToken.Role, s.roles[Director].ID) && !CheckRolePresence(authToken.Role, s.roles[SchoolManager].ID) {
		return []dto.SemesterResponse{}, ErrUnAuthorized
	}

	semesters, err := s.db.GetSemesters(ctx)
	if err != nil {
		return []dto.SemesterResponse{}, err
	}

	for _, s := range semesters {
		semester := dto.SemesterResponse{
			ID:        s.ID,
			Name:      s.Name,
			StartDate: s.StartDate,
			EndDate:   s.EndDate,
			CreatedBy: s.CreatedBy.Int64,
			UpdatedBy: s.UpdatedBy.Int64,
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
		}
		response = append(response, semester)
	}

	return response, err

}

func (s *semesterService) GetCurrentSemester(ctx context.Context, authToken *token.Payload) (dto.SemesterResponse, error) {
	if CheckRolePresence(authToken.Role, s.roles[Admin].ID) {
		return dto.SemesterResponse{}, ErrUnAuthorized
	}
	semester, err := s.db.GetCurrentSemester(ctx)
	if err != nil {
		return dto.SemesterResponse{}, err
	}
	return dto.SemesterResponse{ID: semester.ID,
		Name:      semester.Name,
		StartDate: semester.StartDate,
		EndDate:   semester.EndDate,
		CreatedBy: semester.CreatedBy.Int64,
		UpdatedBy: semester.UpdatedBy.Int64,
		CreatedAt: semester.CreatedAt.Time,
		UpdatedAt: semester.UpdatedAt.Time}, err
}

func NewSemesterService(database db.Store, mapRoles map[string]db.Role) SemesterService {
	return &semesterService{db: database, roles: mapRoles}
}
