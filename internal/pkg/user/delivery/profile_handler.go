package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	userusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
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

	userIDStr := r.URL.Query().Get("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		errMessageParse := fmt.Sprintf("user id = %s a должен быть строкой", userIDStr)
		u.logger.Errorf("%w %s", err, errMessageParse)
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrBadRequest, errMessageParse))

		return
	}

	user, err := u.storage.GetUserWithoutPasswordByID(ctx, userID)
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
		delivery.SendErrResponse(w, u.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	userWithoutPassword := &models.UserWithoutPassword{
		ID: userID,
	}

	if r.Method == http.MethodPatch {
		userWithoutPassword, err = userusecases.ValidatePartOfUserWithoutPassword(r.Body)
		if err != nil {
			u.logger.Errorf("in PartiallyUpdateUserHandler: %+v\n", err)
			delivery.HandleErr(w, u.logger, err)

			return
		}
	} else {
		userWithoutPassword, err = userusecases.ValidateUserWithoutPassword(r.Body)
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
