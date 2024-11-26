package models

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
	Title    *string `json:"title" validation:"required"`
	ColumnID *int64  `json:"columnId"`
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
	Title           *string `json:"title" validate:"min=3,max=50"`
	IsDone          *bool   `json:"isDone"`
	PreviousFieldID *int64  `json:"previousFieldId"`
	NextFieldID     *int64  `json:"nextFieldId"`
}

type CheckListFieldPostRequest struct {
	Title string `json:"title" validate:"required"`
}

type CardMoveRequest struct {
	NewColumnID    *int64 `json:"newColumnId" validate:"required"`
	PreviousCardID *int64 `json:"previousCardId" validate:"required"`
	NextCardID     *int64 `json:"NextCardId" validate:"required"`
}

type ColumnMoveRequest struct {
	PreviousColumnID *int64 `json:"previousColumnId" validate:"required"`
	NextColumnID     *int64 `json:"NextColumnId" validate:"required"`
}

type AssignUserRequest struct {
	NickName string `json:"nickname" validate:"required"`
}
