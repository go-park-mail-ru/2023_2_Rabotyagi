package handler

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
)


// addPostHandler godoc
//
//  @Summary    add post
//  @Description  add post by data
//  @Accept      json
//  @Produce    json
//  @Param      post  path    prePost  true  "Post"
//  @Success    200  Response
//  @Failure    400  string
//  @Router      /post/add [post]
func (h *PostHandler) addPostHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")

	h.storage.AddPost(prePost)

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulAddPost)

	log.Printf("added user: %v", prePost)
}

// getPostHandler godoc
//
//  @Summary    get post
//  @Description  get post by id
//  @Accept      json
//  @Produce    json
//  @Param      id  path    uint64  true  "Post ID"
//  @Success    200  Response
//  @Failure    400  string
//  @Router      /post/get/ [get]
func (h *PostHandler) getPostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	type jsonID struct {
		id uint64
	} 
	
	decoder := json.NewDecoder(r.Body)

	postID := new(jsonID)
	if err := decoder.Decode(postID); err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}

	post, err := h.storage.GetPost(postID.id)
	if err != nil {
		log.Printf("post with this id is not exists %v\n", postID )
		sendResponse(w, ErrPostNotExist)
	
		return
	}

	ResponseSuccessfulGetPost := PostResponse{
		Status: StatusResponseSuccessful, 
		Body: *post,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulGetPost)

	log.Printf("added user: %v", post)
}

// getPostsListHandler godoc
//
//  @Summary    get posts list
//  @Description  get posts by count
//  @Accept      json
//  @Produce    json
//  @Param      count  path    int  true  "Posts count"
//  @Success    200  Response
//  @Failure    400  string
//  @Router      /post/get_list [get]
func (h *PostHandler) getPostsListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)
	}

	type count struct {
		count uint64
	} 
	
	decoder := json.NewDecoder(r.Body)

	postsCount := new(count)
	if err := decoder.Decode(postsCount); err != nil {
		log.Printf("%v\n", err)
		sendResponse(w, ErrBadRequest)

		return
	}
	
	posts, err := h.storage.GetNPosts(int(postsCount.count))
	if err != nil {
		log.Printf("n > posts count %v\n", postsCount.count )
		sendResponse(w, ErrPostNotExist)
	
		return
	}

	ResponseSuccessfulGetPostsList := PostsListResponse{
		Status: StatusResponseSuccessful, 
		Body: *posts,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	sendResponse(w, ResponseSuccessfulGetPostsList)

	log.Printf("added user: %v", posts)
}