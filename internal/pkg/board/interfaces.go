package board

import (
	"RPO_back/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type BoardUsecase interface {
	CreateNewBoard(userID int, data models.CreateBoardRequest) (newBoard *models.Board, err error)
	UpdateBoard(userID int, boardID int, data models.BoardPutRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(userID int, boardID int) error
	GetMyBoards(userID int) (boards []models.Board, err error)
	GetMembersPermissions(userID int, boardID int) (data []models.MemberWithPermissions, err error)
	AddMember(userID int, boardID int, addRequest *models.AddMemberRequest) (newMember *models.MemberWithPermissions, err error)
	UpdateMemberRole(userID int, boardID int, memberID int, newRole string) (updatedMember *models.MemberWithPermissions, err error)
	RemoveMember(userID int, boardID int, memberID int) error
	GetBoardContent(userID int, boardID int) (content *models.BoardContent, err error)
	CreateNewCard(userID int, boardID int, data *models.CardPutRequest) (newCard *models.Card, err error)
	UpdateCard(userID int, boardID int, cardID int, data *models.CardPutRequest) (updatedCard *models.Card, err error)
	DeleteCard(userID int, boardID int, cardID int) (err error)
	CreateColumn(userID int, boardID int, data *models.ColumnRequest) (newCol *models.Column, err error)
	UpdateColumn(userID int, boardID int, columnID int, data *models.ColumnRequest) (updatedCol *models.Column, err error)
	DeleteColumn(userID int, boardID int, columnID int) (err error)
}

type BoardRepo interface {
	CreateBoard(name string, userID int) (*models.Board, error)
	GetBoard(boardID int) (*models.Board, error)
	UpdateBoard(boardID int, data *models.BoardPutRequest) (updatedBoard *models.Board, err error)
	DeleteBoard(boardId int) error
	GetBoardsForUser(userID int) (boardArray []models.Board, err error)
	GetCardsForBoard(boardID int) (cards []models.Card, err error)
	GetColumnsForBoard(boardID int) (columns []models.Column, err error)
	CreateNewCard(boardID int, columnID int, title string) (newCard *models.Card, err error)
	UpdateCard(boardID int, cardID int, data models.CardPutRequest) (updateCard *models.Card, err error)
	DeleteCard(boardID int, cardID int) (err error)
	CreateColumn(boardId int, title string) (newColumn *models.Column, err error)
	UpdateColumn(boardID int, columnID int, data models.ColumnRequest) (updateColumn *models.Column, err error)
	DeleteColumn(boardID int, columnID int) (err error)
	GetUserProfile(userID int) (user *models.UserProfile, err error)
	GetMemberPermissions(boardID int, memberUserID int, getAdderInfo bool) (member *models.MemberWithPermissions, err error)
	GetMembersWithPermissions(boardID int) (members []models.MemberWithPermissions, err error)
	SetMemberRole(boardID int, memberUserID int, newRole string) (member *models.MemberWithPermissions, err error)
	RemoveMember(boardID int, memberUserID int) (err error)
	AddMember(boardID int, adderID int, memberUserID int) (member *models.MemberWithPermissions, err error)
	GetUserByNickname(nickname string) (user *models.UserProfile, err error)
	SetBoardBackground(userID int, boardID int, fileExtension string, fileSize int) (fileName string, err error)
}
