package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	logger, err := mylogger.Get()
	if err != nil {
		http.Error(w, responses.ErrInternalServer, http.StatusInternalServerError)

		return
	}

	responses.SendResponse(w, logger,
		responses.ResponseSuccessful{
			Status: statuses.StatusResponseSuccessful,
			Body:   responses.ResponseBody{Message: "OK"},
		})
}
