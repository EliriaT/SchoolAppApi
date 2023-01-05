// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"
)

type Class struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	HeadTeacher int64         `json:"headTeacher"`
	CreatedBy   sql.NullInt64 `json:"createdBy"`
	UpdatedBy   sql.NullInt64 `json:"updatedBy"`
	CreatedAt   sql.NullTime  `json:"createdAt"`
	UpdatedAt   sql.NullTime  `json:"updatedAt"`
}

type Course struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	TeacherID  int64         `json:"teacherID"`
	SemesterID int64         `json:"semesterID"`
	ClassID    int64         `json:"classID"`
	Dates      []time.Time   `json:"dates"`
	CreatedBy  sql.NullInt64 `json:"createdBy"`
	UpdatedBy  sql.NullInt64 `json:"updatedBy"`
	CreatedAt  sql.NullTime  `json:"createdAt"`
	UpdatedAt  sql.NullTime  `json:"updatedAt"`
}

type Lesson struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	CourseID  int64          `json:"courseID"`
	StartHour time.Time      `json:"startHour"`
	EndHour   time.Time      `json:"endHour"`
	WeekDay   string         `json:"weekDay"`
	Classroom sql.NullString `json:"classroom"`
	CreatedBy sql.NullInt64  `json:"createdBy"`
	UpdatedBy sql.NullInt64  `json:"updatedBy"`
	CreatedAt sql.NullTime   `json:"createdAt"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
}

type Mark struct {
	ID       int64        `json:"id"`
	CourseID int64        `json:"courseID"`
	MarkDate time.Time    `json:"markDate"`
	IsAbsent sql.NullBool `json:"isAbsent"`
	// Bigger than 0, lower than 11
	Mark      sql.NullInt32 `json:"mark"`
	StudentID int64         `json:"studentID"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
	CreatedAt sql.NullTime  `json:"createdAt"`
	UpdatedAt sql.NullTime  `json:"updatedAt"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type School struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
	CreatedAt sql.NullTime  `json:"createdAt"`
	UpdatedAt sql.NullTime  `json:"updatedAt"`
}

type Semester struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	StartDate time.Time     `json:"startDate"`
	EndDate   time.Time     `json:"endDate"`
	CreatedBy sql.NullInt64 `json:"createdBy"`
	UpdatedBy sql.NullInt64 `json:"updatedBy"`
	CreatedAt sql.NullTime  `json:"createdAt"`
	UpdatedAt sql.NullTime  `json:"updatedAt"`
}

type User struct {
	ID                int64          `json:"id"`
	Email             string         `json:"email"`
	Password          string         `json:"password"`
	TotpSecret        string         `json:"totpSecret"`
	LastName          string         `json:"lastName"`
	FirstName         string         `json:"firstName"`
	Gender            string         `json:"gender"`
	PhoneNumber       sql.NullString `json:"phoneNumber"`
	Domicile          sql.NullString `json:"domicile"`
	BirthDate         sql.NullTime   `json:"birthDate"`
	PasswordChangedAt time.Time      `json:"passwordChangedAt"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         sql.NullTime   `json:"updatedAt"`
}

type UserRole struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"userID"`
	RoleID   int64 `json:"roleID"`
	SchoolID int64 `json:"schoolID"`
}

type UserRoleClass struct {
	ID         int64 `json:"id"`
	UserRoleID int64 `json:"userRoleID"`
	ClassID    int64 `json:"classID"`
}
