package usecase

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
	_, err = uc.boardRepository.AddMember(ctx, newBoard.ID, userID, userID)
	if err != nil {
		return nil, err
	}
	_, err = uc.boardRepository.SetMemberRole(ctx, newBoard.ID, userID, "admin")
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
	updatedMember, err = uc.boardRepository.SetMemberRole(ctx, boardID, memberID, newRole)
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

	err = uc.boardRepository.UpdateLastVisit(ctx, userID, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardContent (UpdateLastVisit): %w", err)
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
func (uc *BoardUsecase) UpdateColumn(ctx context.Context, userID int64, boardID int64, columnID int64, data *models.ColumnRequest) (updatedCol *models.Column, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateColumn (get perms): %w", err)
	}
	if perms.Role == "viewer" {
		return nil, fmt.Errorf("UpdateColumn (check): %w", errs.ErrNotPermitted)
	}

	updatedCol, err = uc.boardRepository.UpdateColumn(ctx, boardID, columnID, *data)
	if err != nil {
		return nil, fmt.Errorf("UpdateColumn (add UpdateColumn): %w", err)
	}

	return &models.Column{
		ID:    updatedCol.ID,
		Title: updatedCol.Title,
	}, nil
}

// DeleteColumn удаляет колонку
func (uc *BoardUsecase) DeleteColumn(ctx context.Context, userID int64, boardID int64, columnID int64) (err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return fmt.Errorf("DeleteColumn (get perms): %w", err)
	}
	if perms.Role == "viewer" {
		return fmt.Errorf("DeleteColumn (check): %w", errs.ErrNotPermitted)
	}

	err = uc.boardRepository.DeleteColumn(ctx, boardID, columnID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *BoardUsecase) SetBoardBackground(ctx context.Context, userID int64, boardID int64, file *multipart.File, fileHeader *multipart.FileHeader) (updatedBoard *models.Board, err error) {
	perms, err := uc.boardRepository.GetMemberPermissions(ctx, boardID, userID, false)
	if err != nil {
		return nil, fmt.Errorf("SetBoardBackground (get perms): %w", err)
	}
	if perms.Role != "admin" && perms.Role != "editor_chief" {
		return nil, fmt.Errorf("UpdateColumn (check): %w", errs.ErrNotPermitted)
	}
	uploadTo, err := uc.boardRepository.SetBoardBackground(
		ctx,
		userID,
		boardID,
		uploads.ExtractFileExtension(fileHeader.Filename),
		fileHeader.Size,
	)
	if err != nil {
		return nil, err
	}
	uploadDir := os.Getenv("USER_UPLOADS_DIR")
	filePath := filepath.Join(uploadDir, uploadTo)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("cant create file on server side: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, *file); err != nil {
		return nil, fmt.Errorf("cant copy file on server side: %w", err)
	}
	return uc.boardRepository.GetBoard(ctx, boardID, userID)
}

// AssignUser назначает карточку пользователю
func (uc *BoardUsecase) AssignUser(ctx context.Context, userID int64, cardID int64, assignedUserID int64) (assignedUser *models.UserProfile, err error) {
	panic("not implemented")
}

// DeassignUser отменяет назначение карточки пользователю
func (uc *BoardUsecase) DeassignUser(ctx context.Context, userID int64, cardID int64, assignedUserID int64) (err error) {
	panic("not implemented")
}

// AddComment добавляет комментарий на карточку
func (uc *BoardUsecase) AddComment(ctx context.Context, userID int64, cardID int64, commentReq *models.CommentRequest) (newComment *models.Comment, err error) {
	panic("not implemented")
}

// UpdateComment редактирует существующий комментарий на карточке
func (uc *BoardUsecase) UpdateComment(ctx context.Context, userID int64, commentID int64, commentReq *models.CommentRequest) (updatedComment *models.Comment, err error) {
	panic("not implemented")
}

// DeleteComment удаляет комментарий с карточки
func (uc *BoardUsecase) DeleteComment(ctx context.Context, userID int64, commentID int64) (err error) {
	panic("not implemented")
}

// AddCheckListField добавляет строку чеклиста в конец списка
func (uc *BoardUsecase) AddCheckListField(ctx context.Context, userID int64, cardID int64, fieldReq *models.CheckListFieldPostRequest) (newField *models.CheckListField, err error) {
	panic("not implemented")
}

// UpdateCheckListField обновляет строку чеклиста и/или её положение
func (uc *BoardUsecase) UpdateCheckListField(ctx context.Context, userID int64, fieldID int64, fieldReq *models.CheckListFieldPatchRequest) (updatedField *models.CheckListField, err error) {
	panic("not implemented")
}

// DeleteCheckListField удаляет строку из чеклиста
func (uc *BoardUsecase) DeleteCheckListField(ctx context.Context, userID int64, fieldID int64) (err error) {
	panic("not implemented")
}

// SetCardCover устанавливает обложку для карточки
func (uc *BoardUsecase) SetCardCover(ctx context.Context, userID int64, cardID int64, file []byte, fileHeader *multipart.FileHeader) (updatedCard *models.Card, err error) {
	panic("not implemented")
}

// DeleteCardCover удаляет обложку с карточки
func (uc *BoardUsecase) DeleteCardCover(ctx context.Context, userID int64, cardID int64) (err error) {
	panic("not implemented")
}

// AddAttachment добавляет вложение на карточку
func (uc *BoardUsecase) AddAttachment(ctx context.Context, userID int64, cardID int64, file []byte, fileHeader *multipart.FileHeader) (newAttachment *models.Attachment, err error) {
	panic("not implemented")
}

// DeleteAttachment удаляет вложение с карточки
func (uc *BoardUsecase) DeleteAttachment(ctx context.Context, userID int64, attachmentID int64) (err error) {
	panic("not implemented")
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
	panic("not implemented")
}

// DeleteInviteLink удаляет ссылку-приглашение
func (uc *BoardUsecase) DeleteInviteLink(ctx context.Context, userID int64, boardID int64) (err error) {
	panic("not implemented")
}

// FetchInvite возвращает информацию о приглашении на доску
func (uc *BoardUsecase) FetchInvite(ctx context.Context, inviteUUID string) (board *models.Board, err error) {
	panic("not implemented")
}

// AcceptInvite добавляет пользователя как зрителя на доску
func (uc *BoardUsecase) AcceptInvite(ctx context.Context, userID int64, inviteUUID string) (board *models.Board, err error) {
	panic("not implemented")
}
