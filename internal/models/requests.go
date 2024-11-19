package models

type CardPutRequest struct {
	NewTitle string `json:"title" validate:"required"`
}

type ColumnRequest struct {
	NewTitle string `json:"title" validate:"required"`
}

type AddMemberRequest struct {
	MemberNickname string `json:"nickname"`
}

type UpdateMemberRequest struct {
	NewRole string `json:"newRole"`
}

type BoardPutRequest struct {
	NewName string `json:"name"`
}

type CreateBoardRequest struct {
	Name string `json:"name"`
}
