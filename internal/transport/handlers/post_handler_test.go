package handler_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"io"


	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
	handler "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/handlers"
)

type TestCase struct {
	Post       *storage.PrePost
	Response   []byte
}

func TestAddPostHandler(t *testing.T) {
	prePost := &storage.PrePost{
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

	expectedResponse, _ := json.Marshal(handler.ResponseSuccessfulAddPost) 
	
	var testCase TestCase = TestCase{
		Post:        prePost,
		Response:    expectedResponse,
	}

	
	reqBody, err := json.Marshal(&prePost)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/post/add", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	postStorageMap := storage.NewPostStorageMap()
	postHandler := &handler.PostHandler{
		Storage: postStorageMap,
	}
	
	postHandler.AddPostHandler(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	bodyStr := string(body)
	if bodyStr != string(testCase.Response) {
		t.Errorf("wrong Response: got %+v, expected %+v",
			bodyStr, testCase.Response)
	}
}

func TestGetPostHandler(t *testing.T) {
	prePost := &storage.PrePost{
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

	post := &storage.Post{
		ID: 1,
        AuthorID:        1,
        Title:           "Test Post",
        Description:     "This is a test post",
        Price:           100,
        SafeTransaction: true,
        Delivery:        true,
        City:            "Moscow",
    }

	ResponseSuccessfulGetPost := handler.PostResponse{
		Status: handler.StatusResponseSuccessful,
		Body:   *post,
	}

	expectedResponse, _ := json.Marshal(ResponseSuccessfulGetPost) 
	
	var testCase TestCase = TestCase{
		Post:        prePost,
		Response:    expectedResponse,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/post/get/1", nil)

	w := httptest.NewRecorder()

	postStorageMap := storage.NewPostStorageMap()
	postHandler := &handler.PostHandler{
		Storage: postStorageMap,
	}

	postHandler.Storage.AddPost(prePost)
	
	postHandler.GetPostHandler(w, req)
	
	resp := w.Result()

	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)

	bodyStr := string(body)
	if bodyStr != string(testCase.Response) {
		t.Errorf("wrong Response: got %+v, expected %+v",
			bodyStr, testCase.Response)
	}
}



