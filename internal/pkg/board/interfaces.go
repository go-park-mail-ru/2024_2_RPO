package board

import (
	"RPO_back/internal/models"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type BoardUsecase interface {
	CreateNewBoard(ctx context.Context, userID int64, data models.BoardRequest) (newBoard *models.Board, err error)
	UpdateBoard(ctx context.Context, userID int64, boardID int64, data models.BoardRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(ctx context.Context, userID int64, boardID int64) error
	GetMyBoards(ctx context.Context, userID int64) (boards []models.Board, err error)
	GetMembersPermissions(ctx context.Context, userID int64, boardID int64) (data []models.MemberWithPermissions, err error)
	AddMember(ctx context.Context, userID int64, boardID int64, addRequest *models.AddMemberRequest) (newMember *models.MemberWithPermissions, err error)
	UpdateMemberRole(ctx context.Context, userID int64, boardID int64, memberID int64, newRole string) (updatedMember *models.MemberWithPermissions, err error)
	RemoveMember(ctx context.Context, userID int64, boardID int64, memberID int64) error
	GetBoardContent(ctx context.Context, userID int64, boardID int64) (content *models.BoardContent, err error)
	CreateNewCard(ctx context.Context, userID int64, boardID int64, data *models.CardPostRequest) (newCard *models.Card, err error)
	UpdateCard(ctx context.Context, userID int64, cardID int64, data *models.CardPatchRequest) (updatedCard *models.Card, err error)
	DeleteCard(ctx context.Context, userID int64, cardID int64) (err error)
	CreateColumn(ctx context.Context, userID int64, boardID int64, data *models.ColumnRequest) (newCol *models.Column, err error)
	UpdateColumn(ctx context.Context, userID int64, columnID int64, data *models.ColumnRequest) (updatedCol *models.Column, err error)
	DeleteColumn(ctx context.Context, userID int64, columnID int64) (err error)
	SetBoardBackground(ctx context.Context, userID int64, boardID int64, file *models.UploadedFile) (updatedBoard *models.Board, err error)
	AssignUser(ctx context.Context, userID int64, cardID int64, data *models.AssignUserRequest) (assignedUser *models.UserProfile, err error)
	DeassignUser(ctx context.Context, userID int64, cardID int64, assignedUserID int64) (err error)
	AddComment(ctx context.Context, userID int64, cardID int64, commentReq *models.CommentRequest) (newComment *models.Comment, err error)
	UpdateComment(ctx context.Context, userID int64, commentID int64, commentReq *models.CommentRequest) (updatedComment *models.Comment, err error)
	DeleteComment(ctx context.Context, userID int64, commentID int64) (err error)
	AddCheckListField(ctx context.Context, userID int64, cardID int64, fieldReq *models.CheckListFieldPostRequest) (newField *models.CheckListField, err error)
	UpdateCheckListField(ctx context.Context, userID int64, fieldID int64, fieldReq *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error)
	DeleteCheckListField(ctx context.Context, userID int64, fieldID int64) (err error)
	SetCardCover(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (updatedCard *models.Card, err error)
	DeleteCardCover(ctx context.Context, userID int64, cardID int64) (err error)
	AddAttachment(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (newAttachment *models.Attachment, err error)
	DeleteAttachment(ctx context.Context, userID int64, attachmentID int64) (err error)
	MoveCard(ctx context.Context, userID int64, cardID int64, moveReq *models.CardMoveRequest) (err error)
	MoveColumn(ctx context.Context, userID int64, columnID int64, moveReq *models.ColumnMoveRequest) (err error)
	GetSharedCard(ctx context.Context, userID int64, cardUuid string) (found *models.SharedCardFoundResponse, dummy *models.SharedCardDummyResponse, err error)
	RaiseInviteLink(ctx context.Context, userID int64, boardID int64) (inviteLink *models.InviteLink, err error)
	DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error)
	FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error)
	AcceptInvite(ctx context.Context, userID int64, inviteUUID string) (board *models.Board, err error)
	GetCardDetails(ctx context.Context, userID int64, cardID int64) (details *models.CardDetails, err error)
}

type BoardRepo interface {
	CreateBoard(ctx context.Context, name string, userID int64) (*models.Board, error)
	GetBoard(ctx context.Context, boardID int64, userID int64) (*models.Board, error)
	UpdateBoard(ctx context.Context, boardID int64, userID int64, data *models.BoardRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(ctx context.Context, boardID int64) error
	GetBoardsForUser(ctx context.Context, userID int64) (boardArray []models.Board, err error)
	GetCardsForBoard(ctx context.Context, boardID int64) (cards []models.Card, err error)
	GetColumnsForBoard(ctx context.Context, boardID int64) (columns []models.Column, err error)
	CreateNewCard(ctx context.Context, columnID int64, title string) (newCard *models.Card, err error)
	UpdateCard(ctx context.Context, cardID int64, data models.CardPatchRequest) (updateCard *models.Card, err error)
	DeleteCard(ctx context.Context, cardID int64) (err error)
	CreateColumn(ctx context.Context, boardId int64, title string) (newColumn *models.Column, err error)
	UpdateColumn(ctx context.Context, columnID int64, data models.ColumnRequest) (updateColumn *models.Column, err error)
	DeleteColumn(ctx context.Context, columnID int64) (err error)
	GetUserProfile(ctx context.Context, userID int64) (user *models.UserProfile, err error)
	GetMemberPermissions(ctx context.Context, boardID int64, memberUserID int64, getAdderInfo bool) (member *models.MemberWithPermissions, err error)
	GetMembersWithPermissions(ctx context.Context, boardID int64, userID int64) (members []models.MemberWithPermissions, err error)
	SetMemberRole(ctx context.Context, userID int64, boardID int64, memberUserID int64, newRole string) (member *models.MemberWithPermissions, err error)
	RemoveMember(ctx context.Context, boardID int64, memberUserID int64) (err error)
	AddMember(ctx context.Context, boardID int64, adderID int64, memberUserID int64) (member *models.MemberWithPermissions, err error)
	GetUserByNickname(ctx context.Context, nickname string) (user *models.UserProfile, err error)
	SetBoardBackground(ctx context.Context, userID int64, boardID int64, fileID int64) (newBoard *models.Board, err error)
	GetMemberFromCard(ctx context.Context, userID int64, cardID int64) (role string, boardID int64, err error)
	GetMemberFromCheckListField(ctx context.Context, userID int64, fieldID int64) (role string, boardID int64, cardID int64, err error)
	GetMemberFromAttachment(ctx context.Context, userID int64, attachmentID int64) (role string, boardID int64, cardID int64, err error)
	GetMemberFromColumn(ctx context.Context, userID int64, columnID int64) (role string, boardID int64, err error)
	GetMemberFromComment(ctx context.Context, userID int64, commentID int64) (role string, boardID int64, cardID int64, err error)
	GetCardCheckList(ctx context.Context, cardID int64) (checkList []models.CheckListField, err error)
	GetCardAssignedUsers(ctx context.Context, cardID int64) (assignedUsers []models.UserProfile, err error)
	GetCardComments(ctx context.Context, cardID int64) (comments []models.Comment, err error)
	GetCardAttachments(ctx context.Context, cardID int64) (attachments []models.Attachment, err error)
	GetCardsForMove(ctx context.Context, col1ID int64, col2ID *int64) (column1 []models.Card, column2 []models.Card, err error)
	GetColumnsForMove(ctx context.Context, boardID int64) (columns []models.Column, err error)
	RearrangeCards(ctx context.Context, columnID int64, cards []models.Card) (err error)
	RearrangeColumns(ctx context.Context, columns []models.Column) (err error)
	RearrangeCheckList(ctx context.Context, fields []models.CheckListField) (err error)
	AssignUserToCard(ctx context.Context, cardID int64, userID int64) (assignedUser *models.UserProfile, err error)
	DeassignUserFromCard(ctx context.Context, cardID int64, assignedUserID int64) (err error)
	CreateComment(ctx context.Context, userID int64, cardID int64, comment *models.CommentRequest) (newComment *models.Comment, err error)
	UpdateComment(ctx context.Context, commentID int64, update *models.CommentRequest) (updatedComment *models.Comment, err error)
	DeleteComment(ctx context.Context, commentID int64) (err error)
	CreateCheckListField(ctx context.Context, cardID int64, field *models.CheckListFieldPostRequest) (newField *models.CheckListField, err error)
	UpdateCheckListField(ctx context.Context, fieldID int64, update *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error)
	DeleteCheckListField(ctx context.Context, fieldID int64) error
	SetCardCover(ctx context.Context, userID int64, cardID int64, fileID int64) (updatedCard *models.Card, err error)
	RemoveCardCover(ctx context.Context, cardID int64) (err error)
	AddAttachment(ctx context.Context, userID int64, cardID int64, fileID int64) (newAttachment *models.Attachment, err error)
	RemoveAttachment(ctx context.Context, attachmentID int64) (err error)
	PullInviteLink(ctx context.Context, userID int64, boardID int64) (link *models.InviteLink, err error)
	DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error)
	FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error)
	AcceptInvite(ctx context.Context, userID int64, boardID int64, invitedUserID int64, inviteUUID string) (board *models.Board, err error)
	DeduplicateFile(ctx context.Context, file *models.UploadedFile) (fileNames []string, fileIDs []int64, err error)
	RegisterFile(ctx context.Context, file *models.UploadedFile) (fileID int64, fileUUID string, err error)
}
