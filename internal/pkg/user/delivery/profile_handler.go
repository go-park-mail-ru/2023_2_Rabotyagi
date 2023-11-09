package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

//	 GetUserHandler godoc
//
//		@Summary    get profile
//		@Description  get profile by id
//
// @Tags profile
//
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "user id"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get/{id} [get]
func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userIDStr := delivery.GetPathParam(r.URL.Path)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		u.logger.Errorf("in GetUserHandler: %+v", err)
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := u.storage.GetUserWithoutPasswordByID(ctx, uint64(userID))
	if err != nil {
		u.logger.Errorf("in GetUserHandler: %+v", err)
		delivery.HandleErr(w, u.logger, err)

		return
	}

	user.Sanitize()

	delivery.SendOkResponse(w, u.logger, NewProfileResponse(delivery.StatusResponseSuccessful, user))
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
//	@Param      user  body models.UserWithoutPassword true  "user data for updating"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/update [patch]
//	@Router      /profile/update [put]
func (u *UserHandler) PartiallyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	var err error

	userID := delivery.GetUserIDFromCookie(r, u.logger)

	userWithoutPassword := &models.UserWithoutPassword{
		ID: userID,
	}

	if r.Method == http.MethodPatch {
		userWithoutPassword, err = userusecases.ValidatePartOfUserWithoutPassword(u.logger, r.Body)
		if err != nil {
			u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
			delivery.HandleErr(w, u.logger, err)

			return
		}
	} else {
		userWithoutPassword, err = userusecases.ValidateUserWithoutPassword(u.logger, r.Body)
		if err != nil {
			u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
			delivery.HandleErr(w, u.logger, err)

			return
		}
	}

	updateDataMap := utils.StructToMap(userWithoutPassword)

	delete(updateDataMap, "ID")

	updatedUser, err := u.storage.UpdateUser(ctx, userID, updateDataMap)
	if err != nil {
		u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
		delivery.HandleErr(w, u.logger, err)

		return
	}

	updatedUser.Sanitize()

	delivery.SendOkResponse(w, u.logger, NewProfileResponse(delivery.StatusResponseSuccessful, updatedUser))
	u.logger.Infof("Successfully updated: %+v", userID)
}
