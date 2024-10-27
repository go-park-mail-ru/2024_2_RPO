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
	OwnerUserId int       `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserPermissions struct {
	CanEdit          bool `json:"canEdit"`
	CanShare         bool `json:"canShare"`
	CanInviteMembers bool `json:"canInviteMembers"`
	IsAdmin          bool `json:"isAdmin"`
}
