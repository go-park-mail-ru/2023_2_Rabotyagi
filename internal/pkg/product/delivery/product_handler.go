package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

type PostHandler struct {
	Storage    usecases.IProductStorage
	AddrOrigin string
}

// AddPostHandler godoc
//
//	@Summary    add product
//	@Description  add product by data
//	@Description Error.status can be:
//	@Description StatusErrBadRequest      = 400
//	@Description  StatusErrInternalServer  = 500
//
//	@Accept      json
//	@Produce    json
//	@Param      product  body models.PreProduct true  "product data for adding"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/add [post]
func (p *PostHandler) AddPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	preProduct, err := usecases.ValidatePreProduct(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in AddPostHandler:", err)

		return
	}

	err = p.Storage.AddProduct(ctx, preProduct)
	if err != nil {
		delivery.HandleErr(w, "in AddPostHandler:", err)

		return
	}

	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulAddPost))
	log.Printf("added product: %+v", preProduct)
}

// GetPostHandler godoc
//
//	@Summary    get product
//	@Description  get product by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "product id"
//	@Success    200  {object} PostResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /product/get/{id} [get]
func (p *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	ctx := r.Context()
	postIDStr := utils.GetPathParam(r.URL.Path)

	userID := usecases.GetUserIDFromCookie(r)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("in GetPostHandler: %+v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s product id == %s But shoud be integer", delivery.ErrBadRequest, postIDStr)))

		return
	}

	post, err := p.Storage.GetProduct(ctx, uint64(postID), userID)
	if err != nil {
		log.Printf("in GetPostHandler: product with this id is not exists %+v\n", postID)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrPostNotExist))

		return
	}

	delivery.SendOkResponse(w, NewPostResponse(delivery.StatusResponseSuccessful, post))
	log.Printf("in GetPostHandler: get product: %+v", post)
}

//
//// TODO product list, у нас лежит размер пачки, с фронта прилетает начиная с какого поста брать
//
//// GetPostsListHandler godoc
////
////	@Summary    get posts
////	@Description  get posts by count
////	@Accept      json
////	@Produce    json
////	@Param      count  query uint64 true  "count posts"
////	@Success    200  {object} PostsListResponse
////	@Failure    405  {string} string
////	@Failure    500  {string} string
////	@Failure    222  {object} delivery.ErrorResponse "Error"
////	@Router      /product/get_list [get]
//func (p *PostHandler) GetPostsListHandler(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//	delivery.SetupCORS(w, p.AddrOrigin)
//
//	if r.Method == http.MethodOptions {
//		return
//	}
//
//	if r.Method != http.MethodGet {
//		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
//	}
//
//	countStr := r.URL.Query().Get("count")
//
//	count, err := strconv.Atoi(countStr)
//	if err != nil {
//		log.Printf("in GetPostsListHandler: %+v\n", err)
//		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
//			fmt.Sprintf("%s count posts == %s But shoud be integer", delivery.ErrBadRequest, countStr)))
//
//		return
//	}
//
//	posts, err := p.Storage.GetNProducts(count)
//	if err != nil {
//		log.Printf("in GetPostsListHandler: n > posts count %+v\n", count)
//		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrNoSuchCountOfPosts))
//
//		return
//	}
//
//	delivery.SendOkResponse(w, NewPostsListResponse(delivery.StatusResponseSuccessful, posts))
//	log.Printf("in GetPostsListHandler: get product list: %+v", posts)
//}
