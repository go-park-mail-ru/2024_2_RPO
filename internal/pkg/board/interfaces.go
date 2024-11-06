package board

import (
	"RPO_back/internal/models"
	"context"
	"mime/multipart"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type BoardUsecase interface {
	CreateNewBoard(ctx context.Context, userID int, data models.CreateBoardRequest) (newBoard *models.Board, err error)
	UpdateBoard(ctx context.Context, userID int, boardID int, data models.BoardPutRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(ctx context.Context, userID int, boardID int) error
	GetMyBoards(ctx context.Context, userID int) (boards []models.Board, err error)
	GetMembersPermissions(ctx context.Context, userID int, boardID int) (data []models.MemberWithPermissions, err error)
	AddMember(ctx context.Context, userID int, boardID int, addRequest *models.AddMemberRequest) (newMember *models.MemberWithPermissions, err error)
	UpdateMemberRole(ctx context.Context, userID int, boardID int, memberID int, newRole string) (updatedMember *models.MemberWithPermissions, err error)
	RemoveMember(ctx context.Context, userID int, boardID int, memberID int) error
	GetBoardContent(ctx context.Context, userID int, boardID int) (content *models.BoardContent, err error)
	CreateNewCard(ctx context.Context, userID int, boardID int, data *models.CardPutRequest) (newCard *models.Card, err error)
	UpdateCard(ctx context.Context, userID int, boardID int, cardID int, data *models.CardPutRequest) (updatedCard *models.Card, err error)
	DeleteCard(ctx context.Context, userID int, boardID int, cardID int) (err error)
	CreateColumn(ctx context.Context, userID int, boardID int, data *models.ColumnRequest) (newCol *models.Column, err error)
	UpdateColumn(ctx context.Context, userID int, boardID int, columnID int, data *models.ColumnRequest) (updatedCol *models.Column, err error)
	DeleteColumn(ctx context.Context, userID int, boardID int, columnID int) (err error)
	SetBoardBackground(ctx context.Context, userID int, boardID int, file *multipart.File, fileHeader *multipart.FileHeader) (updatedBoard *models.Board, err error)
}

type BoardRepo interface {
	CreateBoard(ctx context.Context, name string, userID int) (*models.Board, error)
	GetBoard(ctx context.Context, boardID int) (*models.Board, error)
	UpdateBoard(ctx context.Context, boardID int, data *models.BoardPutRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(ctx context.Context, boardId int) error
	GetBoardsForUser(ctx context.Context, userID int) (boardArray []models.Board, err error)
	GetCardsForBoard(ctx context.Context, boardID int) (cards []models.Card, err error)
	GetColumnsForBoard(ctx context.Context, boardID int) (columns []models.Column, err error)
	CreateNewCard(ctx context.Context, boardID int, columnID int, title string) (newCard *models.Card, err error)
	UpdateCard(ctx context.Context, boardID int, cardID int, data models.CardPutRequest) (updateCard *models.Card, err error)
	DeleteCard(ctx context.Context, boardID int, cardID int) (err error)
	CreateColumn(ctx context.Context, boardId int, title string) (newColumn *models.Column, err error)
	UpdateColumn(ctx context.Context, boardID int, columnID int, data models.ColumnRequest) (updateColumn *models.Column, err error)
	DeleteColumn(ctx context.Context, boardID int, columnID int) (err error)
	GetUserProfile(ctx context.Context, userID int) (user *models.UserProfile, err error)
	GetMemberPermissions(ctx context.Context, boardID int, memberUserID int, getAdderInfo bool) (member *models.MemberWithPermissions, err error)
	GetMembersWithPermissions(ctx context.Context, boardID int) (members []models.MemberWithPermissions, err error)
	SetMemberRole(ctx context.Context, boardID int, memberUserID int, newRole string) (member *models.MemberWithPermissions, err error)
	RemoveMember(ctx context.Context, boardID int, memberUserID int) (err error)
	AddMember(ctx context.Context, boardID int, adderID int, memberUserID int) (member *models.MemberWithPermissions, err error)
	GetUserByNickname(ctx context.Context, nickname string) (user *models.UserProfile, err error)
	SetBoardBackground(ctx context.Context, userID int, boardID int, fileExtension string, fileSize int) (fileName string, err error)
}
