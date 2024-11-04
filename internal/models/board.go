package models

import "time"

type CreateBoardRequest struct {
	Name string `json:"name"`
}

type Board struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	BackgroundImageURL string    `json:"backgroundImageUrl,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// MemberWithPermissions - пользователь с правами (в контексте доски)
type MemberWithPermissions struct {
	User              *UserProfile      `json:"user"`
	Role              string            `json:"role"`
	AddedAt           time.Time         `json:"addedAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	AddedBy           *UserProfile      `json:"addedBy"`
	UpdatedBy         *UserProfile      `json:"updatedBy"`
}


type BoardPutRequest struct {
	NewName        string `json:"name"`
	NewDescription string `json:"description"`
}

type BoardContent struct {
	MyRole    string   `json:"myRole"`
	Cards     []Card   `json:"allCards"`
	Columns   []Column `json:"allColumns"`
	BoardInfo *Board   `json:"boardInfo"`
}

type AddMemberRequest struct {
	MemberNickname string `json:"nickname"`
}

type UpdateMemberRequest struct {
	NewRole string `json:"newRole"`
}
