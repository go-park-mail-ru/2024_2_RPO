package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"
)

var roleLevels = map[string]int{
	"viewer":       0,
	"editor":       1,
	"editor_chief": 2,
	"admin":        3,
}

type BoardUsecase struct {
	boardRepository board.BoardRepo
}

func CreateBoardUsecase(boardRepository board.BoardRepo) *BoardUsecase {
	return &BoardUsecase{
		boardRepository: boardRepository,
	}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (uc *BoardUsecase) CreateNewBoard(ctx context.Context, userID int64, data models.BoardRequest) (newBoard *models.Board, err error) {
	newBoard, err = uc.boardRepository.CreateBoard(ctx, data.NewName, userID)
	if err != nil {
		return nil, err
	}

	return newBoard, nil
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (uc *BoardUsecase) UpdateBoard(ctx context.Context, userID int64, boardID int64, data models.BoardRequest) (updatedBoard *models.Board, err error) {
	deleterMember, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (getting editor perm-s): %w", err)
	}
	if deleterMember.Role != "admin" && deleterMember.Role != "editor_chief" {
		return nil, fmt.Errorf("GetMembersPermissions (checking): %w", errs.ErrNotPermitted)
	}
	updatedBoard, err = uc.boardRepository.UpdateBoard(ctx, boardID, userID, &data)
	return
}

// DeleteBoard удаляет доску
func (uc *BoardUsecase) DeleteBoard(ctx context.Context, userID int64, boardID int64) error {
	deleterMember, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return fmt.Errorf("GetMembersPermissions (getting deleter perm-s): %w", err)
	}
	if deleterMember.Role != "admin" {
		return fmt.Errorf("GetMembersPermissions (checking): %w", errs.ErrNotPermitted)
	}
	err = uc.boardRepository.DeleteBoard(ctx, boardID)
	if err != nil {
		return fmt.Errorf("GetMembersPermissions (action): %w", err)
	}
	return nil
}

// GetMyBoards получает все доски для пользователя
func (uc *BoardUsecase) GetMyBoards(ctx context.Context, userID int64) (boards []models.Board, err error) {
	boards, err = uc.boardRepository.GetBoardsForUser(ctx, userID)
	return
}

// GetMembersPermissions получает информацию о ролях всех участников доски
func (uc *BoardUsecase) GetMembersPermissions(ctx context.Context, userID int64, boardID int64) (data []models.MemberWithPermissions, err error) {
	_, err = uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", err)
	}
	permissions, err := uc.boardRepository.GetMembersWithPermissions(ctx, boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (query): %w", err)
	}
	return permissions, nil
}

// AddMember добавляет участника на доску с правами "viewer" и возвращает его права
func (uc *BoardUsecase) AddMember(ctx context.Context, userID int64, boardID int64, addRequest *models.AddMemberRequest) (newMember *models.MemberWithPermissions, err error) {
	adderMember, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", err)
	}
	if (adderMember.Role != "admin") && (adderMember.Role != "editor_chief") {
		return nil, fmt.Errorf("GetMembersPermissions (permissions): %w", errs.ErrNotPermitted)
	}
	newMemberProfile, err := uc.boardRepository.GetUserByNickname(ctx, addRequest.MemberNickname)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (get new user ID): %w", err)
	}
	newMember, err = uc.boardRepository.AddMember(ctx, boardID, userID, newMemberProfile.ID)
	if err != nil {
		return nil, fmt.Errorf("GetMembersPermissions (action): %w", err)
	}
	return newMember, nil
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (uc *BoardUsecase) UpdateMemberRole(ctx context.Context, userID int64, boardID int64, memberID int64, newRole string) (updatedMember *models.MemberWithPermissions, err error) {
	updaterMember, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (updater permissions): %w", err)
	}

	memberToUpdate, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, memberID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (member permissions): %w", err)
	}

	if updaterMember.Role != "admin" {
		if (updaterMember.Role != "admin") && (updaterMember.Role != "editor_chief") {
			return nil, fmt.Errorf("UpdateMemberRole (check1): %w", errs.ErrNotPermitted)
		}
		if roleLevels[updaterMember.Role] <= roleLevels[newRole] {
			return nil, fmt.Errorf("UpdateMemberRole (check2): %w", errs.ErrNotPermitted)
		}
		if roleLevels[updaterMember.Role] <= roleLevels[memberToUpdate.Role] {
			return nil, fmt.Errorf("UpdateMemberRole (check3): %w", errs.ErrNotPermitted)
		}
	}

	updatedMember, err = uc.boardRepository.SetMemberRole(ctx, userID, boardID, memberID, newRole)
	if err != nil {
		return nil, fmt.Errorf("UpdateMemberRole (action): %w", err)
	}

	return updatedMember, nil
}

// RemoveMember удаляет участника с доски
func (uc *BoardUsecase) RemoveMember(ctx context.Context, userID int64, boardID int64, memberID int64) error {
	removerMember, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return fmt.Errorf("RemoveMember (remover permissions): %w", err)
	}
	memberToRemove, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, memberID, false)
	if err != nil {
		return fmt.Errorf("RemoveMember (member permissions): %w", err)
	}
	fmt.Printf("%s removes %s\n", removerMember.Role, memberToRemove.Role)
	if removerMember.Role != "admin" && userID != memberID {
		if (removerMember.Role != "admin") && (removerMember.Role != "editor_chief") {
			return fmt.Errorf("RemoveMember (check1): %w", errs.ErrNotPermitted)
		}

		if roleLevels[removerMember.Role] <= roleLevels[memberToRemove.Role] {
			return fmt.Errorf("RemoveMember (check2): %w", errs.ErrNotPermitted)
		}
	}
	err = uc.boardRepository.RemoveMember(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("UpdateMemberRole (action): %w", err)
	}
	return nil
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (uc *BoardUsecase) GetBoardContent(ctx context.Context, userID int64, boardID int64) (content *models.BoardContent, err error) {
	userPermissions, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("GetBoardContent: %w", errs.ErrNotPermitted)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("GetBoardContent: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetBoardContent (add GetMemberPermissions): %w", err)
	}

	cards, err := uc.boardRepository.GetCardsForBoard(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetCardsForBoard): %w", err)
	}

	cols, err := uc.boardRepository.GetColumnsForBoard(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetColumnsForBoard): %w", err)
	}

	info, err := uc.boardRepository.GetBoard(ctx, boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (add GetBoard): %w", err)
	}

	return &models.BoardContent{
		Cards:     cards,
		Columns:   cols,
		BoardInfo: info,
		MyRole:    userPermissions.Role,
	}, nil
}

// CreateNewCard создаёт новую карточку и возвращает её
func (uc *BoardUsecase) CreateNewCard(ctx context.Context, userID int64, boardID int64, data *models.CardPostRequest) (newCard *models.Card, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("CreateNewCard (get permissions): %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("CreateNewCard (get permissions): %w", err)
		}
		return nil, fmt.Errorf("CreateNewCard (get permissions): %w", err)
	}

	if perms.Role == "viewer" {
		return nil, fmt.Errorf("CreateNewCard (check): %w", errs.ErrNotPermitted)
	}

	card, err := uc.boardRepository.CreateNewCard(ctx, *data.ColumnID, *data.Title)
	if err != nil {
		return nil, fmt.Errorf("CreateNewCard (create): %w", err)
	}

	return &models.Card{
		ID:        card.ID,
		Title:     card.Title,
		ColumnID:  card.ColumnID,
		CreatedAt: card.CreatedAt,
		UpdatedAt: card.UpdatedAt,
	}, nil
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (uc *BoardUsecase) UpdateCard(ctx context.Context, userID int64, cardID int64, data *models.CardPatchRequest) (updatedCard *models.Card, err error) {
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			return nil, fmt.Errorf("UpdateCard (get permissions): %w", err)
		}
		if errors.Is(err, errs.ErrNotFound) {
			return nil, fmt.Errorf("UpdateCard (get permissions): %w", err)
		}
		return nil, fmt.Errorf("UpdateCard (get permissions): %w", err)
	}
	if role == "viewer" {
		return nil, fmt.Errorf("UpdateCard (check): %w", errs.ErrNotPermitted)
	}

	updatedCard, err = uc.boardRepository.UpdateCard(ctx, cardID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateCard (update): %w", err)
	}

	return &models.Card{
		ID:        updatedCard.ID,
		Title:     updatedCard.Title,
		ColumnID:  updatedCard.ColumnID,
		CreatedAt: updatedCard.CreatedAt,
		UpdatedAt: updatedCard.UpdatedAt,
	}, nil
}

// DeleteCard удаляет карточку
func (uc *BoardUsecase) DeleteCard(ctx context.Context, userID int64, cardID int64) (err error) {
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return err
	}
	if role == "viewer" {
		return fmt.Errorf("DeleteCard (check): %w", errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteCard(ctx, cardID)
	if err != nil {
		return fmt.Errorf("DeleteCard (delete): %w", err)
	}

	return nil
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (uc *BoardUsecase) CreateColumn(ctx context.Context, userID int64, boardID int64, data *models.ColumnRequest) (newCol *models.Column, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("CreateColumn (get role): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("CreateColumn (check): %w", errs.ErrNotPermitted)
	}

	column, err := uc.boardRepository.CreateColumn(ctx, boardID, data.NewTitle)
	if err != nil {
		return nil, fmt.Errorf("CreateColumn (create): %w", err)
	}

	return &models.Column{
		ID:    column.ID,
		Title: column.Title,
	}, nil
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (uc *BoardUsecase) UpdateColumn(ctx context.Context, userID int64, columnID int64, data *models.ColumnRequest) (updatedCol *models.Column, err error) {
	role, _, err := uc.boardRepository.GetMemberFromColumn(ctx, userID, columnID)
	if err != nil {
		return nil, fmt.Errorf("UpdateColumn (get perms): %w", err)
	}

	if role == "viewer" {
		return nil, fmt.Errorf("UpdateColumn (check): %w", errs.ErrNotPermitted)
	}

	updatedCol, err = uc.boardRepository.UpdateColumn(ctx, columnID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateColumn (add UpdateColumn): %w", err)
	}

	return updatedCol, nil
}

// DeleteColumn удаляет колонку
func (uc *BoardUsecase) DeleteColumn(ctx context.Context, userID int64, columnID int64) (err error) {
	role, _, err := uc.boardRepository.GetMemberFromColumn(ctx, userID, columnID)
	if err != nil {
		return fmt.Errorf("DeleteColumn (get perms): %w", err)
	}

	if role == "viewer" {
		return fmt.Errorf("DeleteColumn (check): %w", errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteColumn(ctx, columnID)
	if err != nil {
		return fmt.Errorf("DeleteColumn (delete): %w", errs.ErrNotPermitted)
	}

	return nil
}

func (uc *BoardUsecase) SetBoardBackground(ctx context.Context, userID int64, boardID int64, file *models.UploadedFile) (updatedBoard *models.Board, err error) {
	funcName := "SetBoardBackground"
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("SetBoardBackground (get perms): %w", err)
	}
	if perms.Role != "admin" && perms.Role != "editor_chief" {
		return nil, fmt.Errorf("UpdateColumn (check): %w", errs.ErrNotPermitted)
	}

	fileID, err := uploads.UsecaseUploadFile(ctx, file, uc.boardRepository)
	if err != nil {
		return nil, fmt.Errorf("%s (upload): %w", funcName, err)
	}

	newBoard, err := uc.boardRepository.SetBoardBackground(ctx, userID, boardID, fileID)
	if err != nil {
		return nil, fmt.Errorf("%s (set): %w", funcName, err)
	}

	return newBoard, nil
}

// AssignUser назначает карточку пользователю
func (uc *BoardUsecase) AssignUser(ctx context.Context, userID int64, cardID int64, data *models.AssignUserRequest) (assignedUser *models.UserProfile, err error) {
	funcName := "AssignUser"
	perms, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (get perms): %w", funcName, err)
	}

	if perms == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	assignedUserID, err := uc.boardRepository.GetUserByNickname(ctx, data.NickName)
	if err != nil {
		return nil, fmt.Errorf("%s (check): %w", funcName, err)
	}

	assignedUser, err = uc.boardRepository.AssignUserToCard(ctx, cardID, assignedUserID.ID)
	if err != nil {
		return nil, fmt.Errorf("%s (assign user): %w", funcName, err)
	}

	return assignedUser, nil
}

// DeassignUser отменяет назначение карточки пользователю
func (uc *BoardUsecase) DeassignUser(ctx context.Context, userID int64, cardID int64, assignedUserID int64) (err error) {
	funcName := "DeassignUser"
	perms, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return fmt.Errorf("%s (get perms): %w", funcName, err)
	}

	if perms == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeassignUserFromCard(ctx, cardID, assignedUserID)
	if err != nil {
		return fmt.Errorf("%s (deassign user): %w", funcName, err)
	}

	return nil
}

// AddComment добавляет комментарий на карточку
func (uc *BoardUsecase) AddComment(ctx context.Context, userID int64, cardID int64, commentReq *models.CommentRequest) (newComment *models.Comment, err error) {
	funcName := "AddComment"

	perms, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (get perms): %w", funcName, err)
	}

	if perms == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	newComment, err = uc.boardRepository.CreateComment(ctx, userID, cardID, commentReq)
	if err != nil {
		return nil, fmt.Errorf("%s (add comment): %w", funcName, err)
	}

	return newComment, nil
}

// UpdateComment редактирует существующий комментарий на карточке
func (uc *BoardUsecase) UpdateComment(ctx context.Context, userID int64, commentID int64, commentReq *models.CommentRequest) (updatedComment *models.Comment, err error) {
	funcName := "UpdateComment"
	role, _, _, err := uc.boardRepository.GetMemberFromComment(ctx, userID, commentID)
	if err != nil {
		return nil, fmt.Errorf("%s (get perms): %w", funcName, err)
	}

	if role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}
	updatedComment, err = uc.boardRepository.UpdateComment(ctx, commentID, commentReq)
	if err != nil {
		return nil, fmt.Errorf("%s (update comment): %w", funcName, err)
	}

	return updatedComment, nil
}

// DeleteComment удаляет комментарий с карточки
func (uc *BoardUsecase) DeleteComment(ctx context.Context, userID int64, commentID int64) (err error) {
	funcName := "DeleteComment"
	role, _, _, err := uc.boardRepository.GetMemberFromComment(ctx, userID, commentID)
	if err != nil {
		return fmt.Errorf("%s (get perms): %w", funcName, err)
	}

	if role == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("%s (delete comment): %w", funcName, err)
	}

	return nil
}

// AddCheckListField добавляет строку чеклиста в конец списка
func (uc *BoardUsecase) AddCheckListField(ctx context.Context, userID int64, cardID int64, fieldReq *models.CheckListFieldPostRequest) (newField *models.CheckListField, err error) {
	funcName := "AddCheckListField"
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (member): %w", funcName, err)
	}
	if role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	field, err := uc.boardRepository.CreateCheckListField(ctx, cardID, fieldReq)
	if err != nil {
		return nil, fmt.Errorf("%s (create): %w", funcName, err)
	}
	return field, nil
}

// UpdateCheckListField обновляет строку чеклиста и/или её положение
func (uc *BoardUsecase) UpdateCheckListField(ctx context.Context, userID int64, fieldID int64, fieldReq *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error) {
	funcName := "UpdateCheckListField"
	role, _, _, err := uc.boardRepository.GetMemberFromCheckListField(ctx, userID, fieldID)
	if err != nil {
		return nil, fmt.Errorf("%s (member): %w", funcName, err)
	}
	if role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	field, err := uc.boardRepository.UpdateCheckListField(ctx, fieldID, fieldReq)
	if err != nil {
		return nil, fmt.Errorf("%s (update): %w", funcName, err)
	}
	return field, nil
}

// DeleteCheckListField удаляет строку из чеклиста
func (uc *BoardUsecase) DeleteCheckListField(ctx context.Context, userID int64, fieldID int64) (err error) {
	funcName := "DeleteCheckListField"
	role, _, _, err := uc.boardRepository.GetMemberFromCheckListField(ctx, userID, fieldID)
	if err != nil {
		return fmt.Errorf("%s (member): %w", funcName, err)
	}
	if role == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteCheckListField(ctx, fieldID)
	if err != nil {
		return fmt.Errorf("%s (delete): %w", funcName, err)
	}
	return nil
}

// SetCardCover устанавливает обложку для карточки
func (uc *BoardUsecase) SetCardCover(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (updatedCard *models.Card, err error) {
	funcName := "SetCardCover"
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (member): %w", funcName, err)
	}

	if role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	fileID, err := uploads.UsecaseUploadFile(ctx, file, uc.boardRepository)
	if err != nil {
		return nil, fmt.Errorf("%s (upload): %w", funcName, err)
	}

	updatedCard, err = uc.boardRepository.SetCardCover(ctx, userID, cardID, fileID)
	if err != nil {
		return nil, fmt.Errorf("%s (update): %w", funcName, err)
	}

	return updatedCard, nil
}

// DeleteCardCover удаляет обложку с карточки
func (uc *BoardUsecase) DeleteCardCover(ctx context.Context, userID int64, cardID int64) (err error) {
	funcName := "DeleteCardCover"
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return fmt.Errorf("%s (member): %w", funcName, err)
	}

	if role == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.RemoveCardCover(ctx, cardID)
	if err != nil {
		return fmt.Errorf("%s (delete): %w", funcName, err)
	}

	return nil
}

// AddAttachment добавляет вложение на карточку
func (uc *BoardUsecase) AddAttachment(ctx context.Context, userID int64, cardID int64, file *models.UploadedFile) (newAttachment *models.Attachment, err error) {
	funcName := "AddAttachment"
	role, _, err := uc.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (member): %w", funcName, err)
	}

	if role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	fileID, err := uploads.UsecaseUploadFile(ctx, file, uc.boardRepository)
	if err != nil {
		return nil, fmt.Errorf("%s (upload): %w", funcName, err)
	}

	newAttachment, err = uc.boardRepository.AddAttachment(ctx, userID, cardID, fileID, file.OriginalName)
	if err != nil {
		return nil, fmt.Errorf("%s (update): %w", funcName, err)
	}

	return newAttachment, nil
}

// DeleteAttachment удаляет вложение с карточки
func (uc *BoardUsecase) DeleteAttachment(ctx context.Context, userID int64, attachmentID int64) (err error) {
	funcName := "DeleteAttachment"
	role, _, _, err := uc.boardRepository.GetMemberFromAttachment(ctx, userID, attachmentID)
	if err != nil {
		return fmt.Errorf("%s (member): %w", funcName, err)
	}

	if role == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.RemoveAttachment(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("%s (delete): %w", funcName, err)
	}

	return nil
}

// MoveCard перемещает карточку на доске
func (uc *BoardUsecase) MoveCard(ctx context.Context, userID int64, cardID int64, moveReq *models.CardMoveRequest) (err error) {
	panic("not implemented")
}

// MoveColumn перемещает колонку на доске
func (uc *BoardUsecase) MoveColumn(ctx context.Context, userID int64, columnID int64, moveReq *models.ColumnMoveRequest) (err error) {
	panic("not implemented")
}

// GetSharedCard даёт информацию о карточке, которой поделились по ссылке
func (uc *BoardUsecase) GetSharedCard(ctx context.Context, userID int64, cardUuid string) (found *models.SharedCardFoundResponse, dummy *models.SharedCardDummyResponse, err error) {
	panic("not implemented")
}

// RaiseInviteLink устанавливает ссылку-приглашение на доску
func (uc *BoardUsecase) RaiseInviteLink(ctx context.Context, userID int64, boardID int64) (inviteLink *models.InviteLink, err error) {
	funcName := "RaiseInviteLink"
	member, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("%s (member): %w", funcName, err)
	}

	if member.Role == "viewer" {
		return nil, fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	inviteLink, err = uc.boardRepository.PullInviteLink(ctx, userID, boardID)
	if err != nil {
		return nil, fmt.Errorf("%s (update): %w", funcName, err)
	}

	return inviteLink, nil
}

// DeleteInviteLink удаляет ссылку-приглашение
func (uc *BoardUsecase) DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error) {
	funcName := "DeleteInviteLink"
	member, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return fmt.Errorf("%s (member): %w", funcName, err)
	}

	if member.Role == "viewer" {
		return fmt.Errorf("%s (check): %w", funcName, errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteInviteLink(ctx, userID, boardID)
	if err != nil {
		return fmt.Errorf("%s (delete): %w", funcName, err)
	}

	return nil
}

// FetchInvite возвращает информацию о приглашении на доску
func (uc *BoardUsecase) FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error) {
	funcName := "FetchInvite"
	board, err = uc.boardRepository.FetchInvite(ctx, inviteUUID)
	if err != nil {
		return nil, fmt.Errorf("%s (fetch): %w", funcName, err)
	}

	return board, nil
}

// AcceptInvite добавляет пользователя как зрителя на доску
func (uc *BoardUsecase) AcceptInvite(ctx context.Context, userID int64, inviteUUID string) (board *models.Board, err error) {
	panic("not implemented")
}

// GetCardDetails возвращает подробное содержание карточки
func (d *BoardUsecase) GetCardDetails(ctx context.Context, userID int64, cardID int64) (details *models.CardDetails, err error) {
	funcName := "GetCardDetails"
	_, _, err = d.boardRepository.GetMemberFromCard(ctx, userID, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	assignedUsers, err := d.boardRepository.GetCardAssignedUsers(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (assigned): %w", funcName, err)
	}

	attachments, err := d.boardRepository.GetCardAttachments(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (attachments): %w", funcName, err)
	}

	checkList, err := d.boardRepository.GetCardCheckList(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (checklist): %w", funcName, err)
	}

	comments, err := d.boardRepository.GetCardComments(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("%s (comments): %w", funcName, err)
	}

	//TODO убрать это позорище
	card, err := d.boardRepository.UpdateCard(ctx, cardID, models.CardPatchRequest{})
	if err != nil {
		return nil, fmt.Errorf("%s (card): %w", funcName, err)
	}

	fmt.Printf("%#v\n", card)

	return &models.CardDetails{
		Attachments:   attachments,
		CheckList:     checkList,
		Comments:      comments,
		AssignedUsers: assignedUsers,
		Card:          card,
	}, nil
}
