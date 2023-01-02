package dto

import "time"

type ClassResponse struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	HeadTeacher     int64          `json:"head_teacher,omitempty"`
	HeadTeacherName string         `json:"head_teacher_name,omitempty"`
	CreatedBy       int64          `json:"createdBy,omitempty"`
	UpdatedBy       int64          `json:"updatedBy,omitempty"`
	CreatedAt       time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt,omitempty"`
	Pupils          []UserResponse `json:"pupils,omitempty"`
}

type GetClassRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type CreateClassRequest struct {
	Name        string `json:"name" binding:"required"`
	HeadTeacher int64  `json:"head_teacher" binding:"required"`
}

type ChangeHeadTeacherRequest struct {
	ClassID       int64 `json:"class_id" binding:"required"`
	HeadTeacherID int64 `json:"head_teacher_id" binding:"required"`
}
