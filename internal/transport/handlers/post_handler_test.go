package handler_test

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"io/ioutil"


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

	
	// Преобразование JSON-объекта в байты
	reqBody, err := json.Marshal(&prePost)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Создание POST-запроса с тестовым JSON-объектом
	req, err := http.NewRequest(http.MethodPost, "/api/v1/post/add", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Создание ResponseRecorder для записи ответа сервера
	rec := httptest.NewRecorder()

	postStorageMap := storage.NewPostStorageMap()
	postHandler := &handler.PostHandler{
		Storage: postStorageMap,
	}
	
	postHandler.AddPostHandler(rec, req)

	resp := rec.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	bodyStr := string(body)
	if bodyStr != string(testCase.Response) {
		t.Errorf("wrong Response: got %+v, expected %+v",
			bodyStr, testCase.Response)
	}
}
