package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
	resp "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/responses"
)

// AddPostHandler godoc
//
//	@Summary    add post
//	@Description  add post by data
//	@Accept      json
//	@Produce    json
//	@Param      post  body storage.PrePost true  "post data for adding"
//	@Success    200  {object} responses.Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error"
//	@Router      /post/add [post]
func (h *PostHandler) AddPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(w, h.addrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	prePost := new(storage.PrePost)
	if err := decoder.Decode(prePost); err != nil {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrUserNotExits))

		return
	}

	h.Storage.AddPost(prePost)
	resp.SendOkResponse(w, resp.NewResponse(resp.StatusResponseSuccessful, resp.ResponseSuccessfulAddPost))
	log.Printf("added user: %v", prePost)
}

// GetPostHandler godoc
//
//	@Summary    get post
//	@Description  get post by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "post id"
//	@Success    200  {object} responses.PostResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error"
//	@Router      /post/get/{id} [get]
func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(w, h.addrOrigin)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	postIDStr := getPathParam(r.URL.Path)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("%v\n", err)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest,
			fmt.Sprintf("%s post id == %s But shoud be integer", resp.ErrBadRequest, postIDStr)))

		return
	}

	post, err := h.Storage.GetPost(uint64(postID))
	if err != nil {
		log.Printf("post with this id is not exists %v\n", postID)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrPostNotExist))

		return
	}

	resp.SendOkResponse(w, resp.NewPostResponse(resp.StatusResponseSuccessful, post))
	log.Printf("added user: %v", post)
}

// GetPostsListHandler godoc
//
//	@Summary    get posts
//	@Description  get posts by count
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count posts"
//	@Success    200  {object} responses.PostsListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Error"
//	@Router      /post/get_list [get]
func (h *PostHandler) GetPostsListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp.SetupCORS(w, h.addrOrigin)

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
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest,
			fmt.Sprintf("%s count posts == %s But shoud be integer", resp.ErrBadRequest, countStr)))

		return
	}

	posts, err := h.Storage.GetNPosts(count)
	if err != nil {
		log.Printf("n > posts count %v\n", count)
		resp.SendErrResponse(w, resp.NewErrResponse(resp.StatusErrBadRequest, resp.ErrNoSuchCountOfPosts))

		return
	}

	resp.SendOkResponse(w, resp.NewPostsListResponse(resp.StatusResponseSuccessful, posts))
	log.Printf("added user: %v", posts)
}
