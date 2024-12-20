package models

//go:generate easyjson -all requests.go

import (
	"RPO_back/internal/pkg/utils/validate"
	"time"
)

type CardPatchRequest struct {
	NewTitle    *string    `json:"title"`
	NewDeadline *time.Time `json:"deadline"`
	IsDone      *bool      `json:"isDone"`
}

type CardPostRequest struct {
	Title    *string `json:"title" validate:"required"`
	ColumnID *int64  `json:"columnId" validate:"required"`
}

type ColumnRequest struct {
	NewTitle string `json:"title" validate:"required,min=3,max=30"`
}

type AddMemberRequest struct {
	MemberNickname string `json:"nickname" validate:"required"`
}

type UpdateMemberRequest struct {
	NewRole string `json:"newRole" validate:"required"`
}

type BoardRequest struct {
	NewName string `json:"name" validate:"required,min=3"`
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

func (r *UserRegisterRequest) Validate() error {
	if err := validate.CheckPassword(r.Password); err != nil {
		return err
	}

	if err := validate.CheckUserName(r.Name); err != nil {
		return err
	}

	return nil
}

type UserProfileUpdateRequest struct {
	NewName string `json:"name" validate:"required,min=3,max=30"`
	Email   string `json:"email" validate:"required,email"`
}

type CommentRequest struct {
	Text string `json:"text" validate:"required,min=3,max=1024"`
}

type CheckListFieldPatchRequest struct {
	Title           *string `json:"title" validate:"min=3"`
	IsDone          *bool   `json:"isDone"`
	PreviousFieldID *int64  `json:"previousFieldId"`
	NextFieldID     *int64  `json:"nextFieldId"`
}

type TagRequest struct {
	Text  string `json:"text" validate:"required,min=3,max=30"`
	Color string `json:"color" validate:"required,hexcolor,min=4,max=7"`
}

type CheckListFieldPostRequest struct {
	Title *string `json:"title" validate:"required,min=3"`
}

type CardMoveRequest struct {
	NewColumnID    *int64 `json:"newColumnId" validate:"required"`
	PreviousCardID *int64 `json:"previousCardId" validate:"required"`
	NextCardID     *int64 `json:"nextCardId" validate:"required"`
}

type ColumnMoveRequest struct {
	PreviousColumnID *int64 `json:"previousColumnId" validate:"required"`
	NextColumnID     *int64 `json:"nextColumnId" validate:"required"`
}

type AssignUserRequest struct {
	NickName string `json:"nickname" validate:"required"`
}

type ElasticRequest struct {
	Title string `json:"title" validate:"required"`
}
