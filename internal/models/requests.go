package models

import "time"

type CardPatchRequest struct {
	NewTitle    *string    `json:"title"`
	NewDeadline *time.Time `json:"deadline"`
	IsDone      *bool      `json:"isDone"`
}

type ColumnRequest struct {
	NewTitle string `json:"title" validate:"required"`
}

type AddMemberRequest struct {
	MemberNickname string `json:"nickname" validate:"required"`
}

type UpdateMemberRequest struct {
	NewRole string `json:"newRole" validate:"required"`
}

type BoardRequest struct {
	NewName string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,min=8,max=50"`
	OldPassword string `json:"oldPassword" validate:"required"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type UserProfileUpdateRequest struct {
	NewName string `json:"name" validate:"required,min=3,max=30"`
	Email   string `json:"email" validate:"required,email"`
}

type CommentRequest struct {
	Text string `json:"text" validate:"required,min=3,max=1024"`
}
type CheckListFieldPutRequest struct {
	Title           string `json:"" validate:"required"`
	IsDone          bool   `json:"" validate:"re"`
	PreviousFieldID *int   `json:"" validate:""`
	NextFieldID     *int   `json:"" validate:""`
}

type CheckListFieldPostRequest struct {
	Title string
}
