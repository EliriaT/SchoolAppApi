package dto

import (
	"time"
)

type CreateSchoolRequest struct {
	Name string `json:"name" binding:"required"`
}

type SchoolResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"createdBy"`
	UpdatedBy int64     `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetSchoolRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListSchoolRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListSchoolResponse []SchoolResponse
