package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"RPO_back/internal/pkg/utils/uploads"
	"net/http"
	"slices"

	log "github.com/sirupsen/logrus"
)

type BoardDelivery struct {
	boardUsecase board.BoardUsecase
}

func CreateBoardDelivery(boardUsecase board.BoardUsecase) *BoardDelivery {
	return &BoardDelivery{boardUsecase: boardUsecase}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (d *BoardDelivery) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateNewBoard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	data := models.BoardRequest{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn(funcName, ": ", err)
		return
	}

	newBoard, err := d.boardUsecase.CreateNewBoard(r.Context(), userID, data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoJSONResponse(w, newBoard, http.StatusCreated)
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (d *BoardDelivery) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateBoard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		logging.Warn(r.Context(), err)
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	data := models.BoardRequest{}
	err = requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	newBoard, err := d.boardUsecase.UpdateBoard(r.Context(), userID, boardID, data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoJSONResponse(w, newBoard, http.StatusOK)
}

// DeleteBoard удаляет доску
func (d *BoardDelivery) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "DeleteBoard")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	err = d.boardUsecase.DeleteBoard(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "DeleteBoard")
		return
	}
	responses.DoEmptyOkResponse(w)
}

// GetMyBoards получает все доски для пользователя
func (d *BoardDelivery) GetMyBoards(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "GetMyBoards")
	if !ok {
		return
	}

	myBoards, err := d.boardUsecase.GetMyBoards(r.Context(), userID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMyBoards")
		return
	}
	responses.DoJSONResponse(w, myBoards, http.StatusOK)
}

// GetMembersPermissions получает информацию о ролях всех участников доски
func (d *BoardDelivery) GetMembersPermissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "GetMembersPermissions")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	memberPermissions, err := d.boardUsecase.GetMembersPermissions(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMembersPermissions")
		return
	}
	responses.DoJSONResponse(w, memberPermissions, http.StatusOK)
}

// AddMember добавляет участника на доску с правами "viewer" и возвращает его права
func (d *BoardDelivery) AddMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "AddMember")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	data := models.AddMemberRequest{}
	err = requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	newMember, err := d.boardUsecase.AddMember(r.Context(), userID, boardID, &data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "AddMember")
		return
	}
	responses.DoJSONResponse(w, newMember, 200)
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (d *BoardDelivery) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateMemberRole")
	if !ok {
		return
	}
	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	memberID, err := requests.GetIDFromRequest(r, "userId", "user_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	data := models.UpdateMemberRequest{}
	err = requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	if !slices.Contains([]string{"viewer", "editor", "editor_chief", "admin"}, data.NewRole) {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedMember, err := d.boardUsecase.UpdateMemberRole(r.Context(), userID, boardID, memberID, data.NewRole)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "UpdateMemberRole")
		return
	}
	responses.DoJSONResponse(w, updatedMember, 200)
}

// RemoveMember удаляет участника с доски
func (d *BoardDelivery) RemoveMember(w http.ResponseWriter, r *http.Request) {
	funcName := "RemoveMember"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	memberID, err := requests.GetIDFromRequest(r, "userId", "user_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.RemoveMember(r.Context(), userID, boardID, memberID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoEmptyOkResponse(w)
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (d *BoardDelivery) GetBoardContent(w http.ResponseWriter, r *http.Request) {
	funcName := "GetBoardContent"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	content, err := d.boardUsecase.GetBoardContent(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, content, http.StatusOK)
}

// CreateNewCard создаёт новую карточку и возвращает её
func (d *BoardDelivery) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateNewCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	requestData := &models.CardPostRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	newCard, err := d.boardUsecase.CreateNewCard(r.Context(), userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, newCard, http.StatusCreated)
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (d *BoardDelivery) UpdateCard(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardId", "card_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	requestData := &models.CardPatchRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCard, err := d.boardUsecase.UpdateCard(r.Context(), userID, cardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, updatedCard, http.StatusOK)
}

// DeleteCard удаляет карточку
func (d *BoardDelivery) DeleteCard(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardId", "card_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteCard(r.Context(), userID, cardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (d *BoardDelivery) CreateColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	newColumn, err := d.boardUsecase.CreateColumn(r.Context(), userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, newColumn, http.StatusCreated)
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (d *BoardDelivery) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnId", "column_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn("columnID is invalid")
		return
	}

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCol, err := d.boardUsecase.UpdateColumn(r.Context(), userID, columnID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, updatedCol, http.StatusOK)
}

// DeleteColumn удаляет колонку
func (d *BoardDelivery) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnID", "column_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteColumn(r.Context(), userID, columnID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// SetMyAvatar принимает у пользователя файл новой обложки доски,
// сохраняет его и возвращает обновлённую доску
func (d *BoardDelivery) SetBoardBackground(w http.ResponseWriter, r *http.Request) {
	funcName := "SetBoardBackground"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	file, err := uploads.FormFile(r)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "no file found")
		return
	}

	updatedBoard, err := d.boardUsecase.SetBoardBackground(r.Context(), userID, boardID, file)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, updatedBoard, 200)
}

// AssignUser назначает карточку пользователю
func (d *BoardDelivery) AssignUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")

}

// DeassignUser отменяет назначение карточки пользователю
func (d *BoardDelivery) DeassignUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// AddComment добавляет комментарий на карточку
func (d *BoardDelivery) AddComment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// UpdateComment редактирует существующий комментарий на карточке
func (d *BoardDelivery) UpdateComment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// DeleteComment удаляет комментарий с карточки
func (d *BoardDelivery) DeleteComment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// AddCheckListField добавляет строку чеклиста в конец списка
func (d *BoardDelivery) AddCheckListField(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// UpdateCheckListField обновляет строку чеклиста и/или её положение
func (d *BoardDelivery) UpdateCheckListField(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// DeleteCheckListField удаляет строку из чеклиста
func (d *BoardDelivery) DeleteCheckListField(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// SetCardCover устанавливает обложку для карточки
func (d *BoardDelivery) SetCardCover(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// DeleteCardCover удаляет обложку с карточки
func (d *BoardDelivery) DeleteCardCover(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// AddAttachment добавляет вложение на карточку
func (d *BoardDelivery) AddAttachment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// DeleteAttachment удаляет вложение с карточки
func (d *BoardDelivery) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// MoveCard перемещает карточку на доске
func (d *BoardDelivery) MoveCard(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// MoveColumn перемещает колонку на доске
func (d *BoardDelivery) MoveColumn(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// GetSharedCard даёт информацию о карточке, которой поделились по ссылке
func (d *BoardDelivery) GetSharedCard(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// RaiseInviteLink устанавливает ссылку-приглашение на доску
func (d *BoardDelivery) RaiseInviteLink(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// DeleteInviteLink удаляет ссылку-приглашение
func (d *BoardDelivery) DeleteInviteLink(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// FetchInvite возвращает информацию о приглашении на доску
func (d *BoardDelivery) FetchInvite(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// AcceptInvite добавляет пользователя как зрителя на доску
func (d *BoardDelivery) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
