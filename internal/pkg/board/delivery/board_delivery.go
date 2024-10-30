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
	"strconv"
)

type BoardDelivery struct {
	boardUsecase *usecase.BoardUsecase
}

func CreateBoardDelivery(boardUsecase *usecase.BoardUsecase) *BoardDelivery {
	return &BoardDelivery{boardUsecase: boardUsecase}
}

// CreateNewBoard создаёт новую доску и возвращает информацию о ней
func (d *BoardDelivery) CreateNewBoard(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// UpdateBoard обновляет информацию о доске и возвращает обновлённую информацию
func (d *BoardDelivery) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// DeleteBoard удаляет доску
func (d *BoardDelivery) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// GetMyBoards получает все доски для пользователя
func (d *BoardDelivery) GetMyBoards(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// GetMembersPermissions получает информацию о ролях всех участников доски
func (d *BoardDelivery) GetMembersPermissions(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// AddMember добавляет участника на доску с правами "viewer" и возвращает его права
func (d *BoardDelivery) AddMember(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// UpdateMemberRole обновляет роль участника и возвращает обновлённые права
func (d *BoardDelivery) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// RemoveMember удаляет участника с доски
func (d *BoardDelivery) RemoveMember(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

// GetBoardContent получает все карточки и колонки с доски, а также информацию о доске
func (d *BoardDelivery) GetBoardContent(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	var getBoardContentRequest models.BoardContent
	err = requests.GetRequestData(r, &getBoardContentRequest)
	if err != nil {
		responses.DoBadResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	content, err := d.boardUsecase.GetBoardContent(userIDInt, getBoardContentRequest.BoardInfo.Id)
	if err != nil {
		if errors.Is(err, errs.ErrNotPermitted) {
			responses.DoBadResponse(w, 403, "No rights to act")
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			responses.DoBadResponse(w, 404, "No such element was found")
			return
		}
		responses.DoBadResponse(w, 500, "Internal Server Error")
	}

	responses.DoEmptyOkResponce(w)
}

// CreateNewCard создаёт новую карточку и возвращает её
func (d *BoardDelivery) CreateNewCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}

// UpdateCard обновляет карточку и возвращает обновлённую версию
func (d *BoardDelivery) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}

// DeleteCard удаляет карточку
func (d *BoardDelivery) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}

// CreateColumn создаёт колонку канбана на доске и возвращает её
func (d *BoardDelivery) CreateColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}

// UpdateColumn изменяет колонку и возвращает её обновлённую версию
func (d *BoardDelivery) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}

// DeleteColumn удаляет колонку
func (d *BoardDelivery) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	userID, hasUserID := session.UserIDFromContext(r.Context())
	if hasUserID == false {
		responses.DoBadResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
}
