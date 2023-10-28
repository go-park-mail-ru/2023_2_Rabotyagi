package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/post/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/utils"
)

type PostHandler struct {
	Storage    *repository.PostStorageMap
	AddrOrigin string
}

// AddPostHandler godoc
//
//	@Summary    add post
//	@Description  add post by data
//	@Accept      json
//	@Produce    json
//	@Param      post  body models.PrePost true  "post data for adding"
//	@Success    200  {object} delivery.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /post/add [post]
func (p *PostHandler) AddPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	prePost := new(models.PrePost)
	if err := decoder.Decode(prePost); err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, delivery.ErrBadRequest))

		return
	}

	p.Storage.AddPost(prePost)
	delivery.SendOkResponse(w, delivery.NewResponse(delivery.StatusResponseSuccessful, ResponseSuccessfulAddPost))
	log.Printf("added post: %+v", prePost)
}

// GetPostHandler godoc
//
//	@Summary    get post
//	@Description  get post by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "post id"
//	@Success    200  {object} PostResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /post/get/{id} [get]
func (p *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	postIDStr := utils.GetPathParam(r.URL.Path)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s post id == %s But shoud be integer", delivery.ErrBadRequest, postIDStr)))

		return
	}

	post, err := p.Storage.GetPost(uint64(postID))
	if err != nil {
		log.Printf("post with this id is not exists %v\n", postID)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrPostNotExist))

		return
	}

	delivery.SendOkResponse(w, NewPostResponse(delivery.StatusResponseSuccessful, post))
	log.Printf("get post: %+v", post)
}

// GetPostsListHandler godoc
//
//	@Summary    get posts
//	@Description  get posts by count
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count posts"
//	@Success    200  {object} PostsListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /post/get_list [get]
func (p *PostHandler) GetPostsListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	delivery.SetupCORS(w, p.AddrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Printf("%v\n", err)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest,
			fmt.Sprintf("%s count posts == %s But shoud be integer", delivery.ErrBadRequest, countStr)))

		return
	}

	posts, err := p.Storage.GetNPosts(count)
	if err != nil {
		log.Printf("n > posts count %v\n", count)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrNoSuchCountOfPosts))

		return
	}

	delivery.SendOkResponse(w, NewPostsListResponse(delivery.StatusResponseSuccessful, posts))
	log.Printf("get post list: %+v", posts)
}
