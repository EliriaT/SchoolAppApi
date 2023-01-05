package service

import (
	"context"
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/api/token"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"log"
	"time"
)

type LessonService interface {
	CreateLesson(ctx context.Context, authToken *token.Payload, req dto.CreateLessonRequest) (dto.LessonResponse, error)
	// output schedule
	GetLessons(ctx context.Context, authToken *token.Payload) (response dto.ScheduleResponse, err error)

	ChangeLesson(ctx context.Context, authToken *token.Payload, req dto.UpdateLessonParams) (dto.LessonResponse, error)

	// showing a course ID schedule
	GetCourseLessons(ctx context.Context, authToken *token.Payload, req dto.GetCourseLessonsRequest) ([]dto.LessonResponse, error)
}

type lessonService struct {
	db    db.Store
	roles map[string]db.Role
}

// When a lesson is created, course's date are updated, if some error occured further, the lesson should be deleted, changed reverted
func (c *lessonService) CreateLesson(ctx context.Context, authToken *token.Payload, req dto.CreateLessonRequest) (dto.LessonResponse, error) {
	if !CheckRolePresence(authToken.Role, c.roles[Director].ID) && !CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) {
		return dto.LessonResponse{}, ErrUnAuthorized
	}
	lesson, err := c.db.CreateLesson(ctx, db.CreateLessonParams{Name: req.Name, CourseID: req.CourseID, StartHour: req.StartHour, EndHour: req.EndHour, WeekDay: req.WeekDay, Classroom: sql.NullString{String: req.Classroom, Valid: true}})
	if err != nil {
		return dto.LessonResponse{}, err
	}

	var startLessonDate time.Time
	currentSemester, err := c.db.GetCurrentSemester(ctx)
	if err != nil {
		return dto.LessonResponse{}, err
	}
	startSemester := currentSemester.StartDate
	endSemester := currentSemester.EndDate
	endWeek := startSemester.AddDate(0, 0, 7)
	for d := startSemester; d.After(endWeek) == false; d = d.AddDate(0, 0, 1) {
		if d.Weekday().String() == lesson.WeekDay {
			startLessonDate = d
			break
		}
	}

	var courseDates []time.Time

	course, err := c.db.GetCourseByID(ctx, lesson.CourseID)
	if err != nil {
		return dto.LessonResponse{}, err
	}
	log.Println(course.Dates)
	courseDates = course.Dates

	for d := startLessonDate; d.After(endSemester) == false; d = d.AddDate(0, 0, 7) {
		courseDates = append(courseDates, d)
	}

	updatedCourse, err := c.db.UpdateCourseDates(ctx, db.UpdateCourseDatesParams{ID: lesson.CourseID, Dates: courseDates})
	if err != nil {
		return dto.LessonResponse{}, err
	}
	log.Println(updatedCourse.Dates)

	return dto.NewLessonsResponse(lesson), nil
}

// returns the person's schedule, only stuff that has lessons and students can see schedule, and headteacher can see their class schedule
func (c *lessonService) GetLessons(ctx context.Context, authToken *token.Payload) (response dto.ScheduleResponse, err error) {
	if CheckRolePresence(authToken.Role, c.roles[Admin].ID) {
		return response, ErrUnAuthorized
	}
	if CheckRolePresence(authToken.Role, c.roles[Student].ID) {
		schedule, err := c.db.GetClassSchedule(ctx, authToken.ClassID)
		if err != nil {
			return response, err
		}
		for _, lesson := range schedule {
			respLesson := dto.NewLessonsResponse(db.Lesson{lesson.ID, lesson.Name_2, lesson.CourseID, lesson.StartHour, lesson.EndHour, lesson.WeekDay, lesson.Classroom, lesson.CreatedBy, lesson.UpdatedBy, lesson.CreatedAt, lesson.UpdatedAt})
			respLesson.TeacherID = lesson.TeacherID
			respLesson.ClassID = lesson.ClassID
			response.MySchedule = append(response.MySchedule, respLesson)
		}
		return response, err
	}

	userRoles, err := c.db.GetUserRoleByUserId(ctx, authToken.UserID)
	if err != nil {
		return response, err
	}
	for _, ur := range userRoles {
		schedule, err := c.db.GetTeacherSchedule(ctx, ur.ID)
		if err != nil {
			return response, err
		}
		for _, lesson := range schedule {
			respLesson := dto.NewLessonsResponse(db.Lesson{lesson.ID, lesson.Name_2, lesson.CourseID, lesson.StartHour, lesson.EndHour, lesson.WeekDay, lesson.Classroom, lesson.CreatedBy, lesson.UpdatedBy, lesson.CreatedAt, lesson.UpdatedAt})
			respLesson.TeacherID = lesson.TeacherID
			respLesson.ClassID = lesson.ClassID
			response.MySchedule = append(response.MySchedule, respLesson)
		}
	}
	if CheckRolePresence(authToken.Role, c.roles[HeadTeacher].ID) {
		schedule, err := c.db.GetClassSchedule(ctx, authToken.ClassID)
		if err != nil {
			return response, err
		}
		for _, lesson := range schedule {
			respLesson := dto.NewLessonsResponse(db.Lesson{lesson.ID, lesson.Name_2, lesson.CourseID, lesson.StartHour, lesson.EndHour, lesson.WeekDay, lesson.Classroom, lesson.CreatedBy, lesson.UpdatedBy, lesson.CreatedAt, lesson.UpdatedAt})
			respLesson.TeacherID = lesson.TeacherID
			respLesson.ClassID = lesson.ClassID
			response.HeadTeacherClassSchedule = append(response.HeadTeacherClassSchedule, respLesson)
		}
	}

	return response, err
}

func (c *lessonService) ChangeLesson(ctx context.Context, authToken *token.Payload, req dto.UpdateLessonParams) (dto.LessonResponse, error) {
	if !CheckRolePresence(authToken.Role, c.roles[Director].ID) && !CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) {
		return dto.LessonResponse{}, ErrUnAuthorized
	}
	lesson, err := c.db.UpdateLesson(ctx, db.UpdateLessonParams{ID: req.LessonID, StartHour: req.StartHour, EndHour: req.EndHour, WeekDay: req.WeekDay, Classroom: sql.NullString{req.Classroom, true}})
	if err != nil {
		return dto.LessonResponse{}, err
	}
	return dto.NewLessonsResponse(lesson), err
}

func (c *lessonService) GetCourseLessons(ctx context.Context, authToken *token.Payload, req dto.GetCourseLessonsRequest) (respone []dto.LessonResponse, err error) {
	if !CheckRolePresence(authToken.Role, c.roles[Director].ID) && !CheckRolePresence(authToken.Role, c.roles[SchoolManager].ID) {
		return []dto.LessonResponse{}, ErrUnAuthorized
	}
	lessons, err := c.db.GetLessonsOfCourse(ctx, req.CourseID)
	if err != nil {
		return []dto.LessonResponse{}, err
	}
	for _, l := range lessons {
		respone = append(respone, dto.NewLessonsResponse(l))
	}
	return respone, err
}

func NewLessonService(database db.Store, mapRoles map[string]db.Role) LessonService {
	return &lessonService{db: database, roles: mapRoles}
}
