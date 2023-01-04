package service

import (
	"context"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
)

type CourseService interface {
	CreateCourse(ctx context.Context, authToken *token.Payload, req dto.CreateCourseRequest) (dto.CourseResponse, error)
	// output depends on role
	GetCourses(ctx context.Context, authToken *token.Payload) (response []dto.CourseResponse, err error)

	ChangeCourse(ctx context.Context, authToken *token.Payload, req dto.UpdateCourseParams) (dto.CourseResponse, error)

	// showing the catalogue with marks
	GetCourseByID(ctx context.Context, authToken *token.Payload, req dto.GetCourseRequest) (dto.GetCourseMarksResponse, error)
}

type courseService struct {
	db    db.Store
	roles map[string]db.Role
}

func (c *courseService) CreateCourse(ctx context.Context, authToken *token.Payload, req dto.CreateCourseRequest) (dto.CourseResponse, error) {
	if !CheckRolePresence(authToken.Role, c.roles[Director].ID) && !CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) {
		return dto.CourseResponse{}, ErrUnAuthorized
	}
	course, err := c.db.CreateCourse(ctx, db.CreateCourseParams{Name: req.Name, TeacherID: req.TeacherID, SemesterID: req.SemesterID, ClassID: req.ClassID})
	if err != nil {
		return dto.CourseResponse{}, err
	}
	return dto.NewCourseResponse(course), nil
}

func (c *courseService) GetCourses(ctx context.Context, authToken *token.Payload) (response []dto.CourseResponse, err error) {
	if CheckRolePresence(authToken.Role, c.roles[Admin].ID) {
		return []dto.CourseResponse{}, ErrUnAuthorized
	}

	if CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) || CheckRolePresence(authToken.Role, c.roles[Director].ID) {
		courses, err := c.db.GetCoursesOfSchool(ctx, authToken.SchoolID)
		if err != nil {
			return []dto.CourseResponse{}, err
		}
		for _, c := range courses {
			course := dto.NewCourseResponse(db.Course{ID: c.ID, Name: c.Name, TeacherID: c.TeacherID, SemesterID: c.SemesterID, ClassID: c.ClassID, Dates: c.Dates, CreatedBy: c.CreatedBy, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt, UpdatedBy: c.UpdatedBy})
			response = append(response, course)
		}

	}
	if CheckRolePresence(authToken.Role, c.roles[Teacher].ID) {
		userRoles, err := c.db.GetUserRoleByUserId(ctx, authToken.UserID)
		if err != nil {
			return []dto.CourseResponse{}, err
		}
		var userRoleID int64
		for _, ur := range userRoles {
			if ur.RoleID == c.roles[Teacher].ID {
				userRoleID = ur.ID
				break
			}
		}
		courses, err := c.db.ListCoursesOfTeacher(ctx, userRoleID)
		if err != nil {
			return []dto.CourseResponse{}, err
		}

		for _, c := range courses {
			course := dto.NewCourseResponse(c)
			response = append(response, course)
		}

	}

	if CheckRolePresence(authToken.Role, c.roles[HeadTeacher].ID) {

		courses, err := c.db.ListCoursesOfClass(ctx, authToken.ClassID)
		if err != nil {
			return []dto.CourseResponse{}, err
		}

		for _, c := range courses {
			course := dto.NewCourseResponse(c)
			response = append(response, course)
		}
	}

	if CheckRolePresence(authToken.Role, c.roles[Student].ID) {

		courses, err := c.db.ListCoursesOfClass(ctx, authToken.ClassID)
		if err != nil {
			return []dto.CourseResponse{}, err
		}

		for _, c := range courses {
			course := dto.NewCourseResponse(c)
			response = append(response, course)
		}
	}
	return response, err

}

func (c *courseService) ChangeCourse(ctx context.Context, authToken *token.Payload, req dto.UpdateCourseParams) (dto.CourseResponse, error) {
	if !CheckRolePresence(authToken.Role, c.roles[Director].ID) && !CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) {
		return dto.CourseResponse{}, ErrUnAuthorized
	}
	updatedCourse, err := c.db.UpdateCourse(ctx, db.UpdateCourseParams{Name: req.Name, TeacherID: req.TeacherID, SemesterID: req.SemesterID, ClassID: req.ClassID})
	if err != nil {
		return dto.CourseResponse{}, err
	}
	response := dto.NewCourseResponse(updatedCourse)
	return response, err
}

func (c *courseService) GetCourseByID(ctx context.Context, authToken *token.Payload, req dto.GetCourseRequest) (dto.GetCourseMarksResponse, error) {
	if CheckRolePresence(authToken.Role, c.roles[Admin].ID) {
		return dto.GetCourseMarksResponse{}, ErrUnAuthorized
	}
	courseMarks, err := c.db.GetCourseMarks(ctx, req.CourseID)
	if err != nil {
		return dto.GetCourseMarksResponse{}, err
	}

	//we check that the teacher is teaching this subject
	if CheckRolePresence(authToken.Role, c.roles[Teacher].ID) {
		userRoles, err := c.db.GetUserRoleByUserId(ctx, authToken.UserID)
		if err != nil {
			return dto.GetCourseMarksResponse{}, err
		}
		var userRoleID int64
		for _, ur := range userRoles {
			if ur.RoleID == c.roles[Teacher].ID {
				userRoleID = ur.ID
				break
			}
		}
		if courseMarks[0].TeacherID != userRoleID {
			return dto.GetCourseMarksResponse{}, ErrUnAuthorized
		}
	}

	// we check that the student or head teacher  belongs to the class of the course
	if CheckRolePresence(authToken.Role, c.roles[Student].ID) || CheckRolePresence(authToken.Role, c.roles[HeadTeacher].ID) {
		if courseMarks[0].ClassID != authToken.ClassID {
			return dto.GetCourseMarksResponse{}, ErrUnAuthorized
		}
	}
	var validCourseMarks []db.GetCourseMarksRow
	//student should see only their marks
	if CheckRolePresence(authToken.Role, c.roles[Student].ID) {
		userRoles, err := c.db.GetUserRoleByUserId(ctx, authToken.UserID)
		if err != nil {
			return dto.GetCourseMarksResponse{}, err
		}
		// ASSUMING STUDENTS HAVE ONLY ONE ROLE
		for _, mark := range courseMarks {
			if mark.StudentID == userRoles[0].ID {
				validCourseMarks = append(validCourseMarks, mark)
			}
		}

	} else {
		validCourseMarks = courseMarks
	}

	var response dto.GetCourseMarksResponse
	response.CourseName = validCourseMarks[0].Name
	response.TeacherID = validCourseMarks[0].TeacherID
	response.SemesterID = validCourseMarks[0].SemesterID
	response.ClassID = validCourseMarks[0].ClassID
	response.Dates = validCourseMarks[0].Dates

	for _, mark := range validCourseMarks {
		m := dto.MarkResponse{
			MarkID:    mark.ID,
			CourseID:  mark.CourseID,
			MarkDate:  mark.MarkDate,
			IsAbsent:  mark.IsAbsent,
			Mark:      mark.Mark,
			StudentID: mark.StudentID,
			CreatedBy: mark.CreatedBy,
			UpdatedBy: mark.UpdatedBy,
			CreatedAt: mark.CreatedAt,
			UpdatedAt: mark.UpdatedAt,
		}
		response.Marks = append(response.Marks, m)
	}
	return response, err
}

func NewCourseService(database db.Store, mapRoles map[string]db.Role) CourseService {
	return &courseService{db: database, roles: mapRoles}
}
