package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
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

	data := models.CreateBoardRequest{}
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

	data := models.BoardPutRequest{}
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
	requestData := &models.CardPutRequest{}
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

	requestData := &models.CardPutRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCard, err := d.boardUsecase.UpdateCard(r.Context(), userID, boardID, cardID, requestData)
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

	err = d.boardUsecase.DeleteCard(r.Context(), userID, boardID, cardID)
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

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn("boardID is invalid")
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

	updatedCol, err := d.boardUsecase.UpdateColumn(r.Context(), userID, boardID, columnID, requestData)
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

	err = d.boardUsecase.DeleteColumn(r.Context(), userID, boardID, columnID)
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
	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	// Ограничение размера 10 МБ
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		responses.DoBadResponse(w, 400, "bad request")
		log.Warn(funcName, ": ", err)
		return
	}
	defer file.Close()

	updatedBoard, err := d.boardUsecase.SetBoardBackground(r.Context(), userID, boardID, &file, fileHeader)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, funcName)
		return
	}

	responses.DoJSONResponse(w, updatedBoard, 200)
}
