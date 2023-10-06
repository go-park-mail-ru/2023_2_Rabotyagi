package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/responses"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/storage"
	handler "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/transport/handlers"
)

type TestCase struct {
	Post     *storage.PrePost
	Response []byte
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

	expectedResponse, _ := json.Marshal(responses.NewResponse(responses.StatusResponseSuccessful, responses.ResponseSuccessfulAddPost))

	var testCase TestCase = TestCase{
		Post:     prePost,
		Response: expectedResponse,
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
		ID:              1,
		AuthorID:        1,
		Title:           "Test Post",
		Description:     "This is a test post",
		Price:           100,
		SafeTransaction: true,
		Delivery:        true,
		City:            "Moscow",
	}

	ResponseSuccessfulGetPost := responses.PostResponse{
		Status: responses.StatusResponseSuccessful,
		Body:   post,
	}

	expectedResponse, _ := json.Marshal(ResponseSuccessfulGetPost)

	var testCase TestCase = TestCase{
		Post:     prePost,
		Response: expectedResponse,
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

//nolint:funlen
func TestGetPostsListHandlerSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name             string
		inputParamCount  int
		handler          *handler.PostHandler
		postsForStorage  []storage.PrePost
		expectedResponse responses.PostsListResponse
	}

	testCases := [...]TestCase{
		{
			name:            "test basic work",
			inputParamCount: 1,
			handler:         &handler.PostHandler{Storage: storage.NewPostStorageMap()},
			postsForStorage: []storage.PrePost{{
				AuthorID:        1,
				Title:           "Test Post",
				Image: storage.Image{
					Url: "test_url",
					Alt: "test_alt",
				},
				Description:     "This is a test post",
				Price:           100,
				SafeTransaction: true,
				Delivery:        true,
				City:            "Moscow",
			}},
			expectedResponse: responses.PostsListResponse{
				Status: responses.StatusResponseSuccessful,
				Body: []*storage.PostInFeed{{
					ID:              1,
					Title:           "Test Post",
					Image: storage.Image{
						Url: "test_url",
						Alt: "test_alt",
					},
					Price:           100,
					SafeTransaction: true,
					Delivery:        true,
					City:            "Moscow",
				}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			for _, v := range testCase.postsForStorage {
				testCase.handler.Storage.AddPost(&v)
			}

			req := httptest.NewRequest(http.MethodGet,
				fmt.Sprintf("/api/v1/post/get_list?count=%d", testCase.inputParamCount), nil)

			w := httptest.NewRecorder()

			testCase.handler.GetPostsListHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.PostsListResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			if !reflect.DeepEqual(testCase.expectedResponse, resultResponse) {
				t.Errorf("wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}
