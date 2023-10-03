package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)

// addPostHandler godoc
//
//	@Summary    add post
//	@Description  add post by data
//	@Accept      json
//	@Produce    json
//	@Param      post  body storage.PrePost true  "post data for adding"
//	@Success    200  {object} Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    200  {object} ErrorResponse
//	@Router      /post/add [post]
func (h *PostHandler) AddPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)

	prePost := new(storage.PrePost)
	if err := decoder.Decode(prePost); err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	h.Storage.AddPost(prePost)
	w.Header().Set("Content-Type", "application/json")
	sendResponse(w, ResponseSuccessfulAddPost)
	log.Printf("added user: %v", prePost)
}

// getPostHandler godoc
//
//	@Summary    get post
//	@Description  get post by id
//	@Accept      json
//	@Produce    json
//	@Param      id  path uint64 true  "post id"
//	@Success    200  {object} Response
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    200  {object} ErrorResponse
//	@Router      /post/get/{id} [get]
func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	postIDStr := getPathParam(r.URL.Path)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	post, err := h.Storage.GetPost(uint64(postID))
	if err != nil {
		log.Printf("post with this id is not exists %v\n", postID)
		sendResponse(w, ErrPostNotExist)

		return
	}

	ResponseSuccessfulGetPost := PostResponse{
		Status: StatusResponseSuccessful,
		Body:   *post,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulGetPost)

	log.Printf("added user: %v", post)
}

// getPostsListHandler godoc
//
//	@Summary    get posts
//	@Description  get posts by count
//	@Accept      json
//	@Produce    json
//	@Param      count  query uint64 true  "count posts"
//	@Success    200  {object} PostsListResponse
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    200  {object} ErrorResponse
//	@Router      /post/get_list [get]
func (h *PostHandler) GetPostsListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	countStr := r.URL.Query().Get("count")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	posts, err := h.Storage.GetNPosts(count)
	if err != nil {
		log.Printf("n > posts count %v\n", count)
		sendResponse(w, ErrPostNotExist)

		return
	}

	ResponseSuccessfulGetPostsList := PostsListResponse{
		Status: StatusResponseSuccessful,
		Body:   posts,
	}

	sendResponse(w, ResponseSuccessfulGetPostsList)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("added user: %v", posts)
}
