package models

import "time"

type Card struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	ColumnID           int       `json:"columnId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	BackgroundImageURL string    `json:"backgroundImageUrl,omitempty"`
}

type CardPutRequest struct {
	NewTitle    string `json:"title" validate:"required"`
	NewColumnId int    `json:"columnId" validate:"required"`
}

type Column struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type ColumnRequest struct {
	NewTitle string `json:"title" validate:"required"`
}
