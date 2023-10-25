package delivery

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Storage    *repository.AuthStorageMap
	AddrOrigin string
}

// GetProfileHandler godoc
//
//	@Summary    get profile
//	@Description  get profile by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "user id"
//	@Success    200  {object} PostResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get/{id} [get]
func (p *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	userIDStr := utils.GetPathParam(r.URL.Path)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s post id == %s But shoud be integer", delivery.ErrBadRequest, userIDStr)))

		return
	}

	user, err := p.Storage.GetUser(uint64(userID)) // TODO запросом из бд
	if err != nil {
		log.Printf("post with this id is not exists %v\n", userID)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrUserNotExist))

		return
	}

	delivery.SendOkResponse(w, NewProfileResponse(delivery.StatusResponseSuccessful, user))
	log.Printf("received user: %v", user)
}

// PartiallyUpdateProfileHandler godoc
//
//	@Summary    update profile
//	@Description  update some fields of profile
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "user id"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get/{id} [get]
func (p *ProfileHandler) PartiallyUpdateProfileHandler(w http.ResponseWriter, r *http.Request) {

}

// GetProfileHandler godoc
//
//	@Summary    get profile
//	@Description  get profile by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "user id"
//	@Success    200  {object} PostResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get/{id} [get]
func (p *ProfileHandler) FullyUpdateProfileHandler(w http.ResponseWriter, r *http.Request) {

}
