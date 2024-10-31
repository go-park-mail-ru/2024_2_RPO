package delivery

import (
	"RPO_back/internal/errs"
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"errors"
	"net/http"

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
	userID, ok := requests.GetUserIDOrFail(w, r, "CreateNewBoard")
	if !ok {
		return
	}

	data := models.CreateBoardRequest{}
	err := requests.GetRequestData(r, &data)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		log.Warn("CreateNewBoard: ", err)
		return
	}

	newBoard, err := d.boardUsecase.CreateNewBoard(userID, data)
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "CreateNewBoard")
	}
	responses.DoJSONResponce(w, newBoard, http.StatusCreated)
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (d *BoardDelivery) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateBoard")
	if !ok {
		return
	}
}

// DeleteBoard удаляет доску
func (d *BoardDelivery) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "DeleteBoard")
	if !ok {
		return
	}
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
	}
	if err != nil {
		responses.ResponseErrorAndLog(w, err, "GetMembersPermissions")
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
	data := models.AddMemberRequest{}
	requests.GetRequestData(r, &data)
	newMember, err := d.boardUsecase.AddMember(userID, boardID)
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (d *BoardDelivery) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "UpdateMemberRole")
	if !ok {
		return
	}
}

// RemoveMember удаляет участника с доски
func (d *BoardDelivery) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := requests.GetUserIDOrFail(w, r, "RemoveMember")
	if !ok {
		return
	}
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (d *BoardDelivery) GetBoardContent(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if !hasUserID {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	boardID, err := requests.GetIDFromRequest(r, "boardId", "board_")
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	content, err := d.boardUsecase.GetBoardContent(userID, boardID)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoJSONResponce(w, content, http.StatusOK)
}

// CreateNewCard создаёт новую карточку и возвращает её
func (d *BoardDelivery) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoJSONResponce(w, newCard, http.StatusCreated)
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (d *BoardDelivery) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoJSONResponce(w, updatedCard, http.StatusOK)
}

// DeleteCard удаляет карточку
func (d *BoardDelivery) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoEmptyOkResponce(w)
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (d *BoardDelivery) CreateColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoJSONResponce(w, newColumn, http.StatusCreated)
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (d *BoardDelivery) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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

	requestData := &models.ColumnRequest{}
	err = requests.GetRequestData(r, requestData)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedCol, err := d.boardUsecase.UpdateColumn(userID, boardID, cardID, requestData)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoJSONResponce(w, updatedCol, http.StatusCreated)
}

// DeleteColumn удаляет колонку
func (d *BoardDelivery) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
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
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, http.StatusForbidden, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, http.StatusNotFound, "No such element was found")
			return
		}
		responses.DoBadResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responses.DoEmptyOkResponce(w)
}
