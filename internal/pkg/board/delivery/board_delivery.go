package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"RPO_back/internal/pkg/utils/uploads"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

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
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		log.Warn(funcName, ": ", err)
		return
	}

	newBoard, err := d.boardUsecase.CreateNewBoard(r.Context(), userID, data)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}
	responses.DoJSONResponse(r, w, newBoard, http.StatusCreated)
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (d *BoardDelivery) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateBoard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	data := models.BoardRequest{}
	err = requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	newBoard, err := d.boardUsecase.UpdateBoard(r.Context(), userID, boardID, data)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, newBoard, http.StatusOK)
}

// DeleteBoard удаляет доску
func (d *BoardDelivery) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "DeleteBoard")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteBoard(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, "DeleteBoard")
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
		responses.ResponseErrorAndLog(r, w, err, "GetMyBoards")
		return
	}
	responses.DoJSONResponse(r, w, myBoards, http.StatusOK)
}

// GetMembersPermissions получает информацию о ролях всех участников доски
func (d *BoardDelivery) GetMembersPermissions(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "GetMembersPermissions")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}
	memberPermissions, err := d.boardUsecase.GetMembersPermissions(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, "GetMembersPermissions")
		return
	}
	responses.DoJSONResponse(r, w, memberPermissions, http.StatusOK)
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (d *BoardDelivery) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateMemberRole")
	if !ok {
		return
	}
	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}
	memberID, err := requests.GetIDFromRequest(r, "userID", "user_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}
	data := models.UpdateMemberRequest{}
	err = requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}
	if !slices.Contains([]string{"viewer", "editor", "editor_chief", "admin"}, data.NewRole) {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	updatedMember, err := d.boardUsecase.UpdateMemberRole(r.Context(), userID, boardID, memberID, data.NewRole)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, "UpdateMemberRole")
		return
	}
	responses.DoJSONResponse(r, w, updatedMember, 200)
}

// RemoveMember удаляет участника с доски
func (d *BoardDelivery) RemoveMember(w http.ResponseWriter, r *http.Request) {
	funcName := "RemoveMember"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	memberID, err := requests.GetIDFromRequest(r, "userID", "user_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.RemoveMember(r.Context(), userID, boardID, memberID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
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

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	content, err := d.boardUsecase.GetBoardContent(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, content, http.StatusOK)
}

// CreateNewCard создаёт новую карточку и возвращает её
func (d *BoardDelivery) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateNewCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}
	requestData := &models.CardPostRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	*requestData.Title = strings.TrimSpace(*requestData.Title)

	newCard, err := d.boardUsecase.CreateNewCard(r.Context(), userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, newCard, http.StatusCreated)
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (d *BoardDelivery) UpdateCard(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		logging.Error(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	requestData := &models.CardPatchRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		logging.Error(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	if requestData.NewTitle != nil {
		*requestData.NewTitle = strings.TrimSpace(*requestData.NewTitle)
	}

	updatedCard, err := d.boardUsecase.UpdateCard(r.Context(), userID, cardID, requestData)
	if err != nil {
		logging.Error(r.Context(), err)
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, updatedCard, http.StatusOK)
}

// DeleteCard удаляет карточку
func (d *BoardDelivery) DeleteCard(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteCard(r.Context(), userID, cardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

func (d *BoardDelivery) SearchCards(w http.ResponseWriter, r *http.Request) {
	funcName := "SearchCards"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	searchValue := r.URL.Query().Get("")
	if searchValue == "" {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	cards, err := d.boardUsecase.SearchCards(r.Context(), userID, searchValue)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, cards, http.StatusCreated)
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (d *BoardDelivery) CreateColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	newColumn, err := d.boardUsecase.CreateColumn(r.Context(), userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, newColumn, http.StatusCreated)
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (d *BoardDelivery) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnID", "column_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCol, err := d.boardUsecase.UpdateColumn(r.Context(), userID, columnID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, updatedCol, http.StatusOK)
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
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteColumn(r.Context(), userID, columnID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
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
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	file, err := uploads.FormFile(r)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "no file found")
		return
	}
	if file == nil {
		logging.Error(r.Context(), "file is nil, but no error")
		responses.DoBadResponseAndLog(r, w, http.StatusInternalServerError, "internal error")
		return
	}

	updatedBoard, err := d.boardUsecase.SetBoardBackground(r.Context(), userID, boardID, file)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, updatedBoard, 200)
}

// AssignUser назначает карточку пользователю
func (d *BoardDelivery) AssignUser(w http.ResponseWriter, r *http.Request) {
	funcName := "AssignUser"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		logging.Error(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	data := &models.AssignUserRequest{}
	err = requests.GetRequestData(r, data)
	if err != nil {
		logging.Error(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	assignedUser, err := d.boardUsecase.AssignUser(r.Context(), userID, cardID, data)
	if err != nil {
		logging.Error(r.Context(), err)
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, assignedUser, http.StatusOK)
}

// DeassignUser отменяет назначение карточки пользователю
func (d *BoardDelivery) DeassignUser(w http.ResponseWriter, r *http.Request) {
	funcName := "DeassignUser"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	assignedUserID, err := requests.GetIDFromRequest(r, "userID", "user_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeassignUser(r.Context(), userID, cardID, assignedUserID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// AddComment добавляет комментарий на карточку
func (d *BoardDelivery) AddComment(w http.ResponseWriter, r *http.Request) {
	funcName := "AddComment"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	commentReq := &models.CommentRequest{}
	err = json.NewDecoder(r.Body).Decode(commentReq)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	newComment, err := d.boardUsecase.AddComment(r.Context(), userID, cardID, commentReq)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, newComment, http.StatusOK)
}

// UpdateComment редактирует существующий комментарий на карточке
func (d *BoardDelivery) UpdateComment(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateComment"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	commentID, err := requests.GetIDFromRequest(r, "commentID", "comment_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	commentReq := &models.CommentRequest{}
	err = json.NewDecoder(r.Body).Decode(commentReq)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	updatedComment, err := d.boardUsecase.UpdateComment(r.Context(), userID, commentID, commentReq)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, updatedComment, http.StatusOK)
}

// DeleteComment удаляет комментарий с карточки
func (d *BoardDelivery) DeleteComment(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteComment"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	commentID, err := requests.GetIDFromRequest(r, "commentID", "comment_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteComment(r.Context(), userID, commentID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// AddCheckListField добавляет строку чеклиста в конец списка
func (d *BoardDelivery) AddCheckListField(w http.ResponseWriter, r *http.Request) {
	funcName := "AddCheckListField"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	data := &models.CheckListFieldPostRequest{}
	err = requests.GetRequestData(r, data)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	cd, err := d.boardUsecase.AddCheckListField(r.Context(), userID, cardID, data)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, cd, http.StatusCreated)
}

// UpdateCheckListField обновляет строку чеклиста и/или её положение
func (d *BoardDelivery) UpdateCheckListField(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateCheckListField"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	fieldID, err := requests.GetIDFromRequest(r, "fieldID", "field_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	data := &models.CheckListFieldPatchRequest{}
	err = requests.GetRequestData(r, data)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	cd, err := d.boardUsecase.UpdateCheckListField(r.Context(), userID, fieldID, data)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, cd, http.StatusOK)
}

// DeleteCheckListField удаляет строку из чеклиста
func (d *BoardDelivery) DeleteCheckListField(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteCheckListField"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	fieldID, err := requests.GetIDFromRequest(r, "fieldID", "field_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteCheckListField(r.Context(), userID, fieldID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// SetCardCover устанавливает обложку для карточки
func (d *BoardDelivery) SetCardCover(w http.ResponseWriter, r *http.Request) {
	funcName := "SetCardCover"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	file, err := uploads.FormFile(r)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "no file found")
		return
	}

	updatedCard, err := d.boardUsecase.SetCardCover(r.Context(), userID, cardID, file)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, updatedCard, http.StatusOK)
}

// DeleteCardCover удаляет обложку с карточки
func (d *BoardDelivery) DeleteCardCover(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteCardCover"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteCardCover(r.Context(), userID, cardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// AddAttachment добавляет вложение на карточку
func (d *BoardDelivery) AddAttachment(w http.ResponseWriter, r *http.Request) {
	funcName := "AddAttachment"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	file, err := uploads.FormFile(r)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "no file found")
		return
	}

	newAttachment, err := d.boardUsecase.AddAttachment(r.Context(), userID, cardID, file)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, newAttachment, http.StatusCreated)
}

// DeleteAttachment удаляет вложение с карточки
func (d *BoardDelivery) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteAttachment"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	attachmentID, err := requests.GetIDFromRequest(r, "attachmentID", "attachment_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteAttachment(r.Context(), userID, attachmentID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// MoveCard перемещает карточку на доске
func (d *BoardDelivery) MoveCard(w http.ResponseWriter, r *http.Request) {
	funcName := "MoveCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	moveReq := &models.CardMoveRequest{}
	err = json.NewDecoder(r.Body).Decode(moveReq)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.MoveCard(r.Context(), userID, cardID, moveReq)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// MoveColumn перемещает колонку на доске
func (d *BoardDelivery) MoveColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "MoveColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnID", "column_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	moveReq := &models.ColumnMoveRequest{}
	err = json.NewDecoder(r.Body).Decode(moveReq)
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.MoveColumn(r.Context(), userID, columnID, moveReq)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// GetSharedCard даёт информацию о карточке, которой поделились по ссылке
func (d *BoardDelivery) GetSharedCard(w http.ResponseWriter, r *http.Request) {
	funcName := "GetSharedCard"
	userID, ok := session.UserIDFromContext(r.Context())
	if !ok {
		userID = -1
	}

	cardUUID, err := requests.GetUUIDFromRequest(r, "cardUUID")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	found, dummy, err := d.boardUsecase.GetSharedCard(r.Context(), userID, cardUUID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	if found == nil {
		responses.DoJSONResponse(r, w, dummy, http.StatusOK)
		return
	}

	responses.DoJSONResponse(r, w, found, http.StatusOK)
}

// RaiseInviteLink устанавливает ссылку-приглашение на доску
func (d *BoardDelivery) RaiseInviteLink(w http.ResponseWriter, r *http.Request) {
	funcName := "RaiseInviteLink"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	inviteLink, err := d.boardUsecase.RaiseInviteLink(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, inviteLink, http.StatusOK)
}

// DeleteInviteLink удаляет ссылку-приглашение
func (d *BoardDelivery) DeleteInviteLink(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteInviteLink"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardID", "board_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteInviteLink(r.Context(), userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoEmptyOkResponse(w)
}

// FetchInvite возвращает информацию о приглашении на доску
func (d *BoardDelivery) FetchInvite(w http.ResponseWriter, r *http.Request) {
	funcName := "FetchInvite"
	inviteUUID, err := requests.GetUUIDFromRequest(r, "inviteUUID")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	board, err := d.boardUsecase.FetchInvite(r.Context(), inviteUUID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, board, http.StatusOK)
}

// AcceptInvite добавляет пользователя как зрителя на доску
func (d *BoardDelivery) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	funcName := "AcceptInvite"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	inviteUUID, err := requests.GetUUIDFromRequest(r, "inviteUUID")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	board, err := d.boardUsecase.AcceptInvite(r.Context(), userID, inviteUUID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, board, http.StatusOK)
}

// GetCardDetails возвращает подробное содержание карточки
func (d *BoardDelivery) GetCardDetails(w http.ResponseWriter, r *http.Request) {
	funcName := "GetCardDetails"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardID", "card_")
	if err != nil {
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	cd, err := d.boardUsecase.GetCardDetails(r.Context(), userID, cardID)
	if err != nil {
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, cd, http.StatusOK)
}
