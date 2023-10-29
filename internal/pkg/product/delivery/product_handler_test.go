package delivery_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"reflect"
//	"testing"
//
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
//	postdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/delivery"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/product/repository"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
//)
//
//type TestCase struct {
//	Post     *models.PreProduct
//	Response []byte
//}
//
//func TestAddPostHandler(t *testing.T) {
//	prePost := &models.PreProduct{
//		AuthorID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeTransaction: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	expectedResponse, err := json.Marshal(delivery.NewResponse(
//		delivery.StatusResponseSuccessful, postdelivery.ResponseSuccessfulAddPost))
//	if err != nil {
//		t.Fatalf("Failed to marshall expepectedResponse. Error: %v", err)
//	}
//
//	testCase := TestCase{
//		Post:     prePost,
//		Response: expectedResponse,
//	}
//
//	reqBody, err := json.Marshal(&prePost)
//	if err != nil {
//		t.Fatalf("Failed to marshal request body: %v", err)
//	}
//
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/product/add", bytes.NewBuffer(reqBody))
//	if req == nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//
//	w := httptest.NewRecorder()
//
//	postStorageMap := repository.NewPostStorageMap()
//	postHandler := &postdelivery.PostHandler{
//		Storage: postStorageMap,
//	}
//
//	postHandler.AddPostHandler(w, req)
//
//	resp := w.Result()
//
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatalf("Failed to io.ReadAll(). Error: %v", err)
//	}
//
//	bodyStr := string(body)
//	if bodyStr != string(testCase.Response) {
//		t.Errorf("wrong Response: got %+v, expected %+v",
//			bodyStr, testCase.Response)
//	}
//}
//
//func TestGetPostHandler(t *testing.T) {
//	prePost := &models.PreProduct{
//		AuthorID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeTransaction: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	post := &models.Product{
//		ID:              1,
//		AuthorID:        1,
//		Title:           "Test Product",
//		Description:     "This is a test product",
//		Price:           100,
//		SafeTransaction: true,
//		Delivery:        true,
//		City:            "Moscow",
//	}
//
//	ResponseSuccessfulGetPost := postdelivery.PostResponse{
//		Status: delivery.StatusResponseSuccessful,
//		Body:   post,
//	}
//
//	expectedResponse, err := json.Marshal(ResponseSuccessfulGetPost)
//	if err != nil {
//		t.Fatalf("Failed to marshall expepectedResponse. Error: %v", err)
//	}
//
//	testCase := TestCase{
//		Post:     prePost,
//		Response: expectedResponse,
//	}
//
//	req := httptest.NewRequest(http.MethodGet, "/api/v1/product/get/1", nil)
//
//	w := httptest.NewRecorder()
//
//	postStorageMap := repository.NewPostStorageMap()
//	postHandler := &postdelivery.PostHandler{
//		Storage: postStorageMap,
//	}
//
//	postHandler.Storage.AddPost(prePost)
//
//	postHandler.GetPostHandler(w, req)
//
//	resp := w.Result()
//
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatalf("Failed to ReadAll resp.Body: %v", err)
//	}
//
//	bodyStr := string(body)
//	if bodyStr != string(testCase.Response) {
//		t.Errorf("wrong Response: got %+v, expected %+v",
//			bodyStr, testCase.Response)
//	}
//}
//
////nolint:funlen
//func TestGetPostsListHandlerSuccessful(t *testing.T) {
//	t.Parallel()
//
//	type TestCase struct {
//		name             string
//		inputParamCount  int
//		handler          *postdelivery.PostHandler
//		postsForStorage  []models.PreProduct
//		expectedResponse postdelivery.PostsListResponse
//	}
//
//	testCases := [...]TestCase{
//		{
//			name:            "test basic work",
//			inputParamCount: 1,
//			handler:         &postdelivery.PostHandler{Storage: repository.NewPostStorageMap()},
//			postsForStorage: []models.PreProduct{{
//				AuthorID: 1,
//				Title:    "Test Product",
//				Image: models.Image{
//					URL: "test_url",
//					Alt: "test_alt",
//				},
//				Description:     "This is a test product",
//				Price:           100,
//				SafeTransaction: true,
//				Delivery:        true,
//				City:            "Moscow",
//			}},
//			expectedResponse: postdelivery.PostsListResponse{
//				Status: delivery.StatusResponseSuccessful,
//				Body: []*models.ProductInFeed{{
//					ID:    1,
//					Title: "Test Product",
//					Image: models.Image{
//						URL: "test_url",
//						Alt: "test_alt",
//					},
//					Price:           100,
//					SafeTransaction: true,
//					Delivery:        true,
//					City:            "Moscow",
//				}},
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(testCase.name, func(t *testing.T) {
//			t.Parallel()
//
//			for _, v := range testCase.postsForStorage {
//				testCase.handler.Storage.AddPost(&v)
//			}
//
//			req := httptest.NewRequest(http.MethodGet,
//				fmt.Sprintf("/api/v1/product/get_list?count=%d", testCase.inputParamCount), nil)
//
//			w := httptest.NewRecorder()
//
//			testCase.handler.GetPostsListHandler(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			receivedResponse, err := io.ReadAll(resp.Body)
//			if err != nil {
//				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
//			}
//
//			var resultResponse postdelivery.PostsListResponse
//
//			err = json.Unmarshal(receivedResponse, &resultResponse)
//			if err != nil {
//				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
//			}
//
//			if !reflect.DeepEqual(testCase.expectedResponse, resultResponse) {
//				t.Errorf("wrong Response: got %+v, expected %+v",
//					resultResponse, testCase.expectedResponse)
//			}
//		})
//	}
//}
