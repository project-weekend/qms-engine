package model

import "time"

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=50"`
	Description string `json:"description"`
}

type CreateProjectResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
