package model

import "time"

type Student struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStudentRequest struct {
	NIM      string  `json:"nim" validate:"required,nim"`
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Grade    float64 `json:"grade" validate:"min=0,max=4"`
	IsActive bool    `json:"is_active"`
}

type UpdateStudentRequest struct {
	ID       int     `json:"id" validate:"required,number"`
	NIM      string  `json:"nim" validate:"required,nim"`
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	Grade    float64 `json:"grade" validate:"min=0,max=4"`
	IsActive bool    `json:"is_active"`
}

type PatchStudentRequest struct {
	ID       int      `json:"id"`
	NIM      *string  `json:"nim,omitempty" validate:"omitnil,nim"`
	Name     *string  `json:"name,omitempty" validate:"omitnil,min=3,max=100"`
	Grade    *float64 `json:"grade,omitempty" validate:"omitnil,min=0,max=4"`
	IsActive *bool    `json:"is_active,omitempty"`
}

type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	Sort     string `json:"sort"`
	Order    string `json:"order"`
	IsActive *bool  `json:"is_active"`
}

func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}
