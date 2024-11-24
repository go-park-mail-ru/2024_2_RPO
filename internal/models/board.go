package models

import "time"

type Board struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
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
	ID               int64      `json:"id"`
	UUID             string     `json:"cardUuid"`
	Title            string     `json:"title"`
	CoverImageURL    string     `json:"coverImageUrl"`
	ColumnID         int64      `json:"columnId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Deadine          *time.Time `json:"deadline,omitempty"`
	IsDone           bool       `json:"isDone"`
	HasCheckList     bool       `json:"hasCheckList"`
	HasAttachments   bool       `json:"hasAttachments"`
	HasAssignedUsers bool       `json:"hasAssignedUsers"`
	HasComments      bool       `json:"hasComments"`
	OrderIndex       int64      `json:"-"`
}

type Column struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	OrderIndex int64  `json:"-"`
}

type Comment struct {
	ID        int64        `json:"id"`
	Text      int64        `json:"text"`
	IsEdited  bool         `json:"isEdited"`
	CreatedBy *UserProfile `json:"createdBy"`
	CreatedAt time.Time    `json:"createdAt"`
}

type CheckListField struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createdAt"`
	IsDone     bool      `json:"isDone"`
	OrderIndex int64     `json:"-"`
}

type Attachment struct {
	ID           int64     `json:"id"`
	OriginalName string    `json:"originalName"`
	FileName     string    `json:"fileName"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CardDetails struct {
	Card          *Card            `json:"card"`
	CheckList     []CheckListField `json:"checkList"`
	Attachments   []Attachment     `json:"attachments"`
	Comments      []Comment        `json:"comments"`
	AssignedUsers []UserProfile    `json:"assignedUsers"`
}

type InviteLink struct {
	InviteLinkUUID string `json:"inviteLinkUuid"`
}
