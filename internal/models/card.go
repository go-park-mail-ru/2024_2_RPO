package models

import "time"

type Card struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ColumnId    int       `json:"columnId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CardPatchRequest struct {
	NewTitle       string `json:"title"`
	NewDescription string `json:"description"`
}

type Column struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type ColumnRequest struct {
	NewTitle string `json:"title"`
}
