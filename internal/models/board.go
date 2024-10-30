package models

import "time"

type CreateBoardRequest struct {
	Name string `json:"name"`
}

type Board struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	BackgroundImageUrl string    `json:"background,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// MemberPermissions нужна для внутренней логики
type MemberPermissions struct {
	CanEdit          bool
	CanShare         bool
	CanInviteMembers bool
	IsAdmin          bool
}

// MemberWithPermissions - пользователь с правами (в контексте доски)
type MemberWithPermissions struct {
	User              UserProfile       `json:"user"`
	MemberPermissions MemberPermissions `json:"-"`
	Role              string            `json:"role"`
	AddedAt           time.Time         `json:"addedAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	AddedBy           *UserProfile      `json:"addedBy"`
	UpdatedBy         *UserProfile      `json:"updatedBy"`
}

// SetFlags устанавливает флаги в соответствии с ролью участника
func (p *MemberWithPermissions) SetFlags() {
	switch p.Role {
	case "viewer":
		p.MemberPermissions = MemberPermissions{IsAdmin: false, CanEdit: false, CanShare: false, CanInviteMembers: false}
	case "editor":
		p.MemberPermissions = MemberPermissions{IsAdmin: false, CanEdit: true, CanShare: true, CanInviteMembers: false}
	case "editor_chief":
		p.MemberPermissions = MemberPermissions{IsAdmin: false, CanEdit: true, CanShare: true, CanInviteMembers: true}
	case "admin":
		p.MemberPermissions = MemberPermissions{IsAdmin: true, CanEdit: false, CanShare: false, CanInviteMembers: false}
	}
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
