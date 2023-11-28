package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

//	GetUserHandler godoc
//
// @Summary    get profile
// @Description  get profile by id
//
// @Tags profile
//
//	@Produce    json
//	@Param      id  query uint64 true  "user id"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get [get]
func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userID, err := utils.ParseUint64FromRequest(r, "id")
	if err != nil {
		responses.HandleErr(w, u.logger, err)

		return
	}

	user, err := u.service.GetUserWithoutPasswordByID(ctx, userID)
	if err != nil {
		responses.HandleErr(w, u.logger, err)

		return
	}

	responses.SendResponse(w, u.logger, NewProfileResponse(user))
	u.logger.Infof("in GetUserHandler: get product: %+v", user)
}

// PartiallyUpdateUserHandler godoc
//
//	@Summary    update profile
//	@Description  update some fields of profile
//
// @Tags profile
//
//	@Accept      json
//	@Produce    json
//	@Param      user  body models.UserWithoutPassword false  "полностью опционален"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/update [patch]
//	@Router      /profile/update [put]
func (u *UserHandler) PartiallyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	var err error

	userID, err := delivery.GetUserIDFromCookie(r)
	if err != nil {
		responses.HandleErr(w, u.logger, err)

		return
	}

	var updatedUser *models.UserWithoutPassword
	if r.Method == http.MethodPatch {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, true, userID)
	} else {
		updatedUser, err = u.service.UpdateUser(ctx, r.Body, false, userID)
	}

	if err != nil {
		u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
		responses.HandleErr(w, u.logger, err)

		return
	}

	responses.SendResponse(w, u.logger, NewProfileResponse(updatedUser))
	u.logger.Infof("Successfully updated: %+v", userID)
}
