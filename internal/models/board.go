package models

import "time"

type CreateBoardRequest struct {
	Name string `json:"name"`
}

type Board struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Background  string    `json:"background,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UserPermissions нужна для внутренней логики
type UserPermissions struct {
	CanEdit          bool
	CanShare         bool
	CanInviteMembers bool
	IsAdmin          bool
}

type BoardPutRequest struct {
	NewName        string `json:"name"`
	NewDescription string `json:"description"`
}
