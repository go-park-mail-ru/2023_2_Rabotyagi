package delivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

//	 GetUserHandler godoc
//
//		@Summary    get profile
//		@Description  get profile by id
//		@Accept      json
//		@Produce    json
//		@Param      id  path uint64 true  "user id"
//		@Success    200  {object} ProfileResponse
//		@Failure    405  {string} string
//		@Failure    500  {string} string
//		@Failure    222  {object} delivery.ErrorResponse "Error"
//		@Router      /profile/get/{id} [get]
func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
		log.Printf("in GetUserHandler: %+v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := u.storage.GetUserWithoutPasswordByID(ctx, uint64(userID))
	if err != nil {
		delivery.HandleErr(w, "error in GetUserHandler:", err)

		return
	}

	delivery.SendOkResponse(w, NewProfileResponse(delivery.StatusResponseSuccessful, user))
	log.Printf("in GetPostHandler: get product: %+v", user)
}

// PartiallyUpdateUserHandler godoc
//
//	@Summary    update profile
//	@Description  update some fields of profile
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
	defer r.Body.Close()
	delivery.SetupCORS(w, u.addrOrigin, u.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	var userWithoutPassword *models.UserWithoutPassword

	var err error

	if r.Method == http.MethodPatch {
		userWithoutPassword, err = usecases.ValidatePartOfUserWithoutPassword(r.Body)
		if err != nil {
			delivery.HandleErr(w, "in PartiallyUpdateUserHandler:", err)

			return
		}
	} else {
		userWithoutPassword, err = usecases.ValidateUserWithoutPassword(r.Body)
		if err != nil {
			delivery.HandleErr(w, "in PartiallyUpdateUserHandler:", err)

			return
		}
	}

	updateDataMap := utils.StructToMap(userWithoutPassword)

	userID, ok := updateDataMap["ID"].(uint64)
	if !ok {
		log.Printf("in PartiallyUpdateUserHandler: userID isn`t uint64")
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	delete(updateDataMap, "ID")

	updatedUser, err := u.storage.UpdateUser(ctx, userID, updateDataMap)
	if err != nil {
		delivery.HandleErr(w, "in PartiallyUpdateUserHandler:", err)

		return
	}

	delivery.SendOkResponse(w, NewProfileResponse(delivery.StatusResponseSuccessful, updatedUser))
	log.Printf("Successfully updated: %+v", userID)
}
