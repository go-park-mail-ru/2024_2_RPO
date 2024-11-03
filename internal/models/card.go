package models

import "time"

type Card struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	ColumnID           int       `json:"columnId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	BackgroundImageURL string    `json:"backgroundImageUrl,omitempty"`
}

type CardPatchRequest struct {
	NewTitle       string `json:"title"`
	NewDescription string `json:"description"`
	ColumnId       int    `json:"columnId"`
}

type Column struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type ColumnRequest struct {
	NewTitle string `json:"title"`
}
