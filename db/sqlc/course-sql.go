package db

//import (
//	"database/sql"
//	"github.com/lib/pq"
//	"time"
//)
//
//// Code generated by sqlc. DO NOT EDIT.
//// versions:
////   sqlc v1.15.0
//// source: course.sql
//
//package db
//
//import (
//"context"
//"database/sql"
//"time"
//"github.com/lib/pq"
//)
//
//const createCourse = `-- name: CreateCourse :one
//INSERT INTO "Course"(
//    name,teacher_id,semester_id,class_id
//)VALUES (
//            $1,$2,$3,$4
//        ) RETURNING id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at
//`
//
//type CreateCourseParams struct {
//	Name       string `json:"name"`
//	TeacherID  int64  `json:"teacherID"`
//	SemesterID int64  `json:"semesterID"`
//	ClassID    int64  `json:"classID"`
//}
//
//func (q *Queries) CreateCourse(ctx context.Context, arg CreateCourseParams) (Course, error) {
//	row := q.db.QueryRowContext(ctx, createCourse,
//		arg.Name,
//		arg.TeacherID,
//		arg.SemesterID,
//		arg.ClassID,
//	)
//	var i Course
//	err := row.Scan(
//		&i.ID,
//		&i.Name,
//		&i.TeacherID,
//		&i.SemesterID,
//		&i.ClassID,
//		pq.Array(&i.Dates),
//		&i.CreatedBy,
//		&i.UpdatedBy,
//		&i.CreatedAt,
//		&i.UpdatedAt,
//	)
//	return i, err
//}
//
//const getCourseByID = `-- name: GetCourseByID :one
//SELECT id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at FROM "Course"
//WHERE id = $1
//`
//
//func (q *Queries) GetCourseByID(ctx context.Context, id int64) (Course, error) {
//	var times PgTimeArray
//	row := q.db.QueryRowContext(ctx, getCourseByID, id)
//	var i Course
//	err := row.Scan(
//		&i.ID,
//		&i.Name,
//		&i.TeacherID,
//		&i.SemesterID,
//		&i.ClassID,
//		&times,
//		&i.CreatedBy,
//		&i.UpdatedBy,
//		&i.CreatedAt,
//		&i.UpdatedAt,
//	)
//
//	for _, t:= range times{
//		i.Dates = append(i.Dates,t.Time)
//	}
//	return i, err
//}
//
//const getCoursesOfSchool = `-- name: GetCoursesOfSchool :many
//SELECT "Course".id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at, "UserRoles".id, user_id, role_id, school_id FROM "Course"
//INNER JOIN "UserRoles"
//ON  "Course".teacher_id = "UserRoles".id AND "UserRoles".school_id = $1
//ORDER BY name
//`
//
//type GetCoursesOfSchoolRow struct {
//	ID         int64         `json:"id"`
//	Name       string        `json:"name"`
//	TeacherID  int64         `json:"teacherID"`
//	SemesterID int64         `json:"semesterID"`
//	ClassID    int64         `json:"classID"`
//	Dates      []time.Time   `json:"dates"`
//	CreatedBy  sql.NullInt64 `json:"createdBy"`
//	UpdatedBy  sql.NullInt64 `json:"updatedBy"`
//	CreatedAt  sql.NullTime  `json:"createdAt"`
//	UpdatedAt  sql.NullTime  `json:"updatedAt"`
//	ID_2       int64         `json:"id2"`
//	UserID     int64         `json:"userID"`
//	RoleID     int64         `json:"roleID"`
//	SchoolID   int64         `json:"schoolID"`
//}
//
//func (q *Queries) GetCoursesOfSchool(ctx context.Context, schoolID int64) ([]GetCoursesOfSchoolRow, error) {
//	rows, err := q.db.QueryContext(ctx, getCoursesOfSchool, schoolID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	items := []GetCoursesOfSchoolRow{}
//	for rows.Next() {
//		var i GetCoursesOfSchoolRow
//		if err := rows.Scan(
//			&i.ID,
//			&i.Name,
//			&i.TeacherID,
//			&i.SemesterID,
//			&i.ClassID,
//			pq.Array(&i.Dates),
//			&i.CreatedBy,
//			&i.UpdatedBy,
//			&i.CreatedAt,
//			&i.UpdatedAt,
//			&i.ID_2,
//			&i.UserID,
//			&i.RoleID,
//			&i.SchoolID,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const listCoursesOfClass = `-- name: ListCoursesOfClass :many
//SELECT id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at FROM "Course"
//WHERE class_id = $1
//ORDER BY name
//`
//
//func (q *Queries) ListCoursesOfClass(ctx context.Context, classID int64) ([]Course, error) {
//	rows, err := q.db.QueryContext(ctx, listCoursesOfClass, classID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	items := []Course{}
//	for rows.Next() {
//		var i Course
//		if err := rows.Scan(
//			&i.ID,
//			&i.Name,
//			&i.TeacherID,
//			&i.SemesterID,
//			&i.ClassID,
//			pq.Array(&i.Dates),
//			&i.CreatedBy,
//			&i.UpdatedBy,
//			&i.CreatedAt,
//			&i.UpdatedAt,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const listCoursesOfTeacher = `-- name: ListCoursesOfTeacher :many
//SELECT id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at FROM "Course"
//WHERE teacher_id = $1
//ORDER BY name
//`
//
//func (q *Queries) ListCoursesOfTeacher(ctx context.Context, teacherID int64) ([]Course, error) {
//	rows, err := q.db.QueryContext(ctx, listCoursesOfTeacher, teacherID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	items := []Course{}
//	for rows.Next() {
//		var i Course
//		if err := rows.Scan(
//			&i.ID,
//			&i.Name,
//			&i.TeacherID,
//			&i.SemesterID,
//			&i.ClassID,
//			pq.Array(&i.Dates),
//			&i.CreatedBy,
//			&i.UpdatedBy,
//			&i.CreatedAt,
//			&i.UpdatedAt,
//		); err != nil {
//			return nil, err
//		}
//		items = append(items, i)
//	}
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//	if err := rows.Err(); err != nil {
//		return nil, err
//	}
//	return items, nil
//}
//
//const updateCourse = `-- name: UpdateCourse :one
//UPDATE  "Course"
//SET  teacher_id = $3, name = $2,semester_id=$5,class_id = $4,updated_at = now()
//where id = $1
//RETURNING id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at
//`
//
//type UpdateCourseParams struct {
//	ID         int64  `json:"id"`
//	Name       string `json:"name"`
//	TeacherID  int64  `json:"teacherID"`
//	ClassID    int64  `json:"classID"`
//	SemesterID int64  `json:"semesterID"`
//}
//
//func (q *Queries) UpdateCourse(ctx context.Context, arg UpdateCourseParams) (Course, error) {
//	row := q.db.QueryRowContext(ctx, updateCourse,
//		arg.ID,
//		arg.Name,
//		arg.TeacherID,
//		arg.ClassID,
//		arg.SemesterID,
//	)
//	var i Course
//	err := row.Scan(
//		&i.ID,
//		&i.Name,
//		&i.TeacherID,
//		&i.SemesterID,
//		&i.ClassID,
//		pq.Array(&i.Dates),
//		&i.CreatedBy,
//		&i.UpdatedBy,
//		&i.CreatedAt,
//		&i.UpdatedAt,
//	)
//	return i, err
//}
//
//const updateCourseDates = `-- name: UpdateCourseDates :one
//UPDATE  "Course"
//SET  dates = $2, updated_at = now()
//where id = $1
//RETURNING id, name, teacher_id, semester_id, class_id, dates, created_by, updated_by, created_at, updated_at
//`
//
//type UpdateCourseDatesParams struct {
//	ID    int64       `json:"id"`
//	Dates []time.Time `json:"dates"`
//}
//
//func (q *Queries) UpdateCourseDates(ctx context.Context, arg UpdateCourseDatesParams) (Course, error) {
//	var times PgTimeArray
//	row := q.db.QueryRowContext(ctx, updateCourseDates, arg.ID, pq.Array(arg.Dates))
//	var i Course
//
//	err := row.Scan(
//		&i.ID,
//		&i.Name,
//		&i.TeacherID,
//		&i.SemesterID,
//		&i.ClassID,
//		&times,
//		&i.CreatedBy,
//		&i.UpdatedBy,
//		&i.CreatedAt,
//		&i.UpdatedAt,
//	)
//	for _, t:= range times{
//		i.Dates = append(i.Dates,t.Time)
//	}
//	return i, err
//}
