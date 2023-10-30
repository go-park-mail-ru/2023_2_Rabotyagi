package delivery

import (
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

type ProfileHandler struct {
	Storage    usecases.IUserStorage
	AddrOrigin string
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	value := reflect.ValueOf(data)
	typeOf := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typeOf.Field(i).Name

		if field.IsZero() {
			continue
		}

		result[fieldName] = field.Interface()
	}

	return result
}

// GetUserHandler godoc
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
func (p *ProfileHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userIDStr := utils.GetPathParam(r.URL.Path)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("in GetUserHandler: %+v", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	user, err := p.Storage.GetUserWithoutPasswordByID(ctx, uint64(userID))
	if err != nil {
		handleErr(w, "error in GetUserHandler:", err)

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
//	@Param      id  path uint64 true  "user id"
//	@Success    200  {object} ProfileResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /profile/get/{id} [get]
func (p *ProfileHandler) PartiallyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()
}

// GetUserHandler godoc
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
func (p *ProfileHandler) FullyUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	ctx := r.Context()

	userWithoutPassword, err := usecases.ValidateUserWithoutID(r.Body)
	if err != nil {
		handleErr(w, "in SignUpHandler:", err)

		return
	}

	updataDataMap := structToMap(userWithoutPassword)

	err := p.Storage.UpdateUser(ctx, userWithoutID)
	if err != nil {
		handleErr(w, "error in SignUpHandler:", err)

		return
	}
}
