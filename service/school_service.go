package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

type SchoolService interface {
	CreateSchool(ctx context.Context, authToken *token.Payload, req dto.CreateSchoolRequest) (dto.SchoolResponse, error)
	GetSchoolByID(ctx context.Context, authToken *token.Payload, req dto.GetSchoolRequest) (dto.SchoolResponse, error)
	ListSchools(ctx context.Context, authToken *token.Payload, req dto.ListSchoolRequest) (dto.ListSchoolResponse, error)
}

type schoolService struct {
	db    db.Store
	roles map[string]db.Role
}

func (s *schoolService) CreateSchool(ctx context.Context, authToken *token.Payload, req dto.CreateSchoolRequest) (dto.SchoolResponse, error) {

	if !CheckRolePresence(authToken.Role, s.roles[Admin].ID) {
		return dto.SchoolResponse{}, ErrUnAuthorized
	}

	school, err := s.db.CreateSchool(ctx, req.Name)
	if err != nil {
		return dto.SchoolResponse{}, err
	}

	response := dto.SchoolResponse{
		ID:        school.ID,
		Name:      school.Name,
		CreatedBy: school.CreatedBy.Int64,
		UpdatedBy: school.UpdatedBy.Int64,
		CreatedAt: school.CreatedAt.Time,
		UpdatedAt: school.UpdatedAt.Time,
	}
	return response, err

}

func (s *schoolService) GetSchoolByID(ctx context.Context, authToken *token.Payload, req dto.GetSchoolRequest) (dto.SchoolResponse, error) {
	if !CheckRolePresence(authToken.Role, s.roles[Admin].ID) && !CheckRolePresence(authToken.Role, s.roles[Director].ID) {
		return dto.SchoolResponse{}, ErrUnAuthorized
	}
	if CheckRolePresence(authToken.Role, s.roles[Director].ID) {
		req.ID = authToken.SchoolID
	}

	school, err := s.db.GetSchoolbyId(ctx, req.ID)
	if err != nil {
		return dto.SchoolResponse{}, err
	}
	response := dto.SchoolResponse{
		ID:        school.ID,
		Name:      school.Name,
		CreatedBy: school.CreatedBy.Int64,
		UpdatedBy: school.UpdatedBy.Int64,
		CreatedAt: school.CreatedAt.Time,
		UpdatedAt: school.UpdatedAt.Time,
	}
	return response, err

}

func (s *schoolService) ListSchools(ctx context.Context, authToken *token.Payload, req dto.ListSchoolRequest) (dto.ListSchoolResponse, error) {
	if !CheckRolePresence(authToken.Role, s.roles[Admin].ID) {
		return dto.ListSchoolResponse{}, ErrUnAuthorized
	}

	arg := db.ListSchoolsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	schools, err := s.db.ListSchools(ctx, arg)
	if err != nil {
		return dto.ListSchoolResponse{}, err
	}

	response := dto.ListSchoolResponse{}

	for _, s := range schools {
		school := dto.SchoolResponse{
			ID:        s.ID,
			Name:      s.Name,
			CreatedBy: s.CreatedBy.Int64,
			UpdatedBy: s.UpdatedBy.Int64,
			CreatedAt: s.CreatedAt.Time,
			UpdatedAt: s.UpdatedAt.Time,
		}
		response = append(response, school)
	}
	return response, err

}

func NewSchoolService(database db.Store, mapRoles map[string]db.Role) SchoolService {
	return &schoolService{db: database, roles: mapRoles}
}
