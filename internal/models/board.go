package models

import "time"

type Board struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	BackgroundImageURL string    `json:"backgroundImageUrl"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	LastVisitAt        time.Time `json:"lastVisitAt"`
}

// MemberWithPermissions - пользователь с правами (в контексте доски)
type MemberWithPermissions struct {
	User      *UserProfile `json:"user"`
	Role      string       `json:"role"`
	AddedAt   time.Time    `json:"addedAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	AddedBy   *UserProfile `json:"addedBy"`
	UpdatedBy *UserProfile `json:"updatedBy"`
}

type BoardContent struct {
	MyRole    string   `json:"myRole"`
	Cards     []Card   `json:"allCards"`
	Columns   []Column `json:"allColumns"`
	BoardInfo *Board   `json:"boardInfo"`
}

type Card struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	ColumnID           int       `json:"columnId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	BackgroundImageURL string    `json:"backgroundImageUrl"`
}

type Column struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
