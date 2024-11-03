package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
	"slices"

	log "github.com/sirupsen/logrus"
)

type BoardDelivery struct {
	boardUsecase *usecase.BoardUsecase
}

func CreateBoardDelivery(boardUsecase *usecase.BoardUsecase) *BoardDelivery {
	return &BoardDelivery{boardUsecase: boardUsecase}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (d *BoardDelivery) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	funcName := "CreateNewBoard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	data := models.CreateBoardRequest{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn(funcName, ": ", err)
		return
	}

	newBoard, err := d.boardUsecase.CreateNewBoard(userID, data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoJSONResponce(w, newBoard, http.StatusCreated)
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (d *BoardDelivery) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateBoard")
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	data := models.BoardPutRequest{}
	err = requests.GetRequestData(r, &data)
	if err == nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	newBoard, err := d.boardUsecase.UpdateBoard(userID, boardID, data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "UpdateBoard")
		return
	}
	responses.DoJSONResponce(w, newBoard, http.StatusOK)
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
	err = d.boardUsecase.DeleteBoard(userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "DeleteBoard")
		return
	}
	responses.DoEmptyOkResponce(w)
}

// GetMyBoards получает все доски для пользователя
func (d *BoardDelivery) GetMyBoards(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "GetMyBoards")
	if !ok {
		return
	}

	myBoards, err := d.boardUsecase.GetMyBoards(userID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMyBoards")
		return
	}
	responses.DoJSONResponce(w, myBoards, http.StatusOK)
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
	memberPermissions, err := d.boardUsecase.GetMembersPermissions(userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMembersPermissions")
		return
	}
	responses.DoJSONResponce(w, memberPermissions, http.StatusOK)
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

	newMember, err := d.boardUsecase.AddMember(userID, boardID, &data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "AddMember")
		return
	}
	responses.DoJSONResponce(w, newMember, 200)
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

	updatedMember, err := d.boardUsecase.UpdateMemberRole(userID, boardID, memberID, data.NewRole)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "UpdateMemberRole")
		return
	}
	responses.DoJSONResponce(w, updatedMember, 200)
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

	err = d.boardUsecase.RemoveMember(userID, boardID, memberID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}
	responses.DoEmptyOkResponce(w)
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

	content, err := d.boardUsecase.GetBoardContent(userID, boardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponce(w, content, http.StatusOK)
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
	requestData := &models.CardPatchRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	newCard, err := d.boardUsecase.CreateNewCard(userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponce(w, newCard, http.StatusCreated)
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (d *BoardDelivery) UpdateCard(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
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

	updatedCard, err := d.boardUsecase.UpdateCard(userID, boardID, cardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponce(w, updatedCard, http.StatusOK)
}

// DeleteCard удаляет карточку
func (d *BoardDelivery) DeleteCard(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteCard"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	cardID, err := requests.GetIDFromRequest(r, "cardId", "card_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteCard(userID, boardID, cardID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoEmptyOkResponce(w)
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

	newColumn, err := d.boardUsecase.CreateColumn(userID, boardID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponce(w, newColumn, http.StatusCreated)
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (d *BoardDelivery) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "UpdateColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn("boardID is invalid")
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnID", "column_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn("boardID is invalid")
		return
	}

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCol, err := d.boardUsecase.UpdateColumn(userID, boardID, columnID, requestData)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponce(w, updatedCol, http.StatusCreated)
}

// DeleteColumn удаляет колонку
func (d *BoardDelivery) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	funcName := "DeleteColumn"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	columnID, err := requests.GetIDFromRequest(r, "columnId", "column_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.boardUsecase.DeleteColumn(userID, boardID, columnID)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoEmptyOkResponce(w)
}
