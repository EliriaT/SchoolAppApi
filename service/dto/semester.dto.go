package dto

import "time"

type CreateSemesterRequest struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required,ltefield=EndDate" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" binding:"required" time_format:"2006-01-02"`
}

type SemesterResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedBy int64     `json:"createdBy"`
	UpdatedBy int64     `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetSemesterRequest struct {
	ID int64 `uri:"id" binding:"required"`
}
