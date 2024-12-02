package delivery

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/poll"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/requests"
	"RPO_back/internal/pkg/utils/responses"
	"net/http"
)

type PollDelivery struct {
	pollUC poll.PollUsecase
}

func CreatePollDelivery(pollUC poll.PollUsecase) *PollDelivery {
	return &PollDelivery{pollUC: pollUC}
}

func (d *PollDelivery) SubmitPoll(w http.ResponseWriter, r *http.Request) {
	funcName := "SubmitPoll"
	userID, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}
	pollSubmit := models.PollSubmit{}
	err := requests.GetRequestData(r, &pollSubmit)
	if err != nil {
		logging.Warn(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusBadRequest, "bad request")
		return
	}

	err = d.pollUC.SubmitPoll(r.Context(), userID, &pollSubmit)
	if err != nil {
		logging.Warn(r.Context(), err)
		responses.DoBadResponseAndLog(r, w, http.StatusInternalServerError, "internal error")
		return
	}

	responses.DoEmptyOkResponse(w)
}

func (d *PollDelivery) GetPollResults(w http.ResponseWriter, r *http.Request) {
	funcName := "GetPollResults"
	_, ok := requests.GetUserIDOrFail(w, r, funcName)
	if !ok {
		return
	}

	pollResults, err := d.pollUC.GetPollResults(r.Context())
	if err != nil {
		logging.Warn(r.Context(), err)
		responses.ResponseErrorAndLog(r, w, err, funcName)
		return
	}

	responses.DoJSONResponse(r, w, pollResults, http.StatusOK)
}
