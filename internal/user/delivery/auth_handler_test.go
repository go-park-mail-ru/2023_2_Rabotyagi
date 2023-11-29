package delivery_test

// import (
//	"bytes"
//	"encoding/json"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"reflect"
//	"testing"
//	"time"
//
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/models"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
//	userdelivery "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/delivery"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/user/repository"
//)
//
//func TestSignUpHandlerSuccessful(t *testing.T) {
//	t.Parallel()
//
//	type TestCase struct {
//		name             string
//		inputPreUser     *models.PreUser
//		expectedResponse *delivery.Response
//	}
//
//	testCases := [...]TestCase{
//		{
//			name:             "test basic work",
//			inputPreUser:     &models.PreUser{Email: "example@mail.ru", Password: "password"},
//			expectedResponse: delivery.NewResponse(delivery.StatusResponseSuccessful, userdelivery.ResponseSuccessfulSignUp),
//		},
//	}
//
//	for _, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(testCase.name, func(t *testing.T) {
//			t.Parallel()
//
//			reqBody, err := json.Marshal(&testCase.inputPreUser)
//			if err != nil {
//				t.Fatalf("Failed to marshal request body: %v", err)
//			}
//
//			req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewBuffer(reqBody))
//
//			w := httptest.NewRecorder()
//
//			authStorageMap := repository.NewAuthStorageMap()
//			authHandler := &userdelivery.AuthHandler{
//				storage: authStorageMap,
//			}
//
//			authHandler.SignUpHandler(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			receivedResponse, err := io.ReadAll(resp.Body)
//			if err != nil {
//				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
//			}
//
//			var resultResponse delivery.Response
//
//			err = json.Unmarshal(receivedResponse, &resultResponse)
//			if err != nil {
//				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
//			}
//
//			if !reflect.DeepEqual(*testCase.expectedResponse, resultResponse) {
//				t.Errorf("wrong Response: got %+v, expected %+v",
//					resultResponse, testCase.expectedResponse)
//			}
//		})
//	}
//}
//
////nolint:funlen
//func TestSignInHandlerSuccessful(t *testing.T) {
//	t.Parallel()
//
//	type TestCase struct {
//		name             string
//		inputPreUser     *models.PreUser
//		expectedResponse *delivery.Response
//	}
//
//	testCases := [...]TestCase{
//		{
//			name:             "test basic work",
//			inputPreUser:     &models.PreUser{Email: "example@mail.ru", Password: "password"},
//			expectedResponse: delivery.NewResponse(delivery.StatusResponseSuccessful, userdelivery.ResponseSuccessfulSignIn),
//		},
//	}
//
//	for _, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(testCase.name, func(t *testing.T) {
//			t.Parallel()
//
//			reqBody, err := json.Marshal(&testCase.inputPreUser)
//			if err != nil {
//				t.Fatalf("Failed to marshal request body: %v", err)
//			}
//
//			req := httptest.NewRequest(http.MethodPost, "/api/v1/signin", bytes.NewBuffer(reqBody))
//
//			w := httptest.NewRecorder()
//
//			authStorageMap := repository.NewAuthStorageMap()
//			err = authStorageMap.CreateUser(testCase.inputPreUser)
//			if err != nil {
//				t.Fatalf("Failed to CreateUser err: %v", err)
//			}
//
//			authHandler := &userdelivery.AuthHandler{
//				storage: authStorageMap,
//			}
//
//			authHandler.SignInHandler(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			receivedResponse, err := io.ReadAll(resp.Body)
//			if err != nil {
//				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
//			}
//
//			var resultResponse delivery.Response
//
//			err = json.Unmarshal(receivedResponse, &resultResponse)
//			if err != nil {
//				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
//			}
//
//			if !reflect.DeepEqual(*testCase.expectedResponse, resultResponse) {
//				t.Errorf("wrong Response: got %+v, expected %+v",
//					resultResponse, testCase.expectedResponse)
//			}
//		})
//	}
//}
//
////nolint:funlen
//func TestLogOutHandlerSuccessful(t *testing.T) {
//	t.Parallel()
//
//	type TestCase struct {
//		name             string
//		inputCookie      *http.Cookie
//		expectedResponse *delivery.Response
//	}
//
//	testCases := [...]TestCase{
//		{
//			name: "test basic work",
//			inputCookie: &http.Cookie{
//				Name: userdelivery.CookieAuthName,
//				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
//					"eyJlbWFpbCI6ImV4YW1wbGVAbWFpbC5ydSIsImV4cGlyZSI6MCwidXNlcklEIjoxfQ." +
//					"GBCEb3XJ6aHTsyl8jC3lxSWK6byjbYN0kg2e3NH2i9s",
//				Expires: time.Now().Add(time.Hour)},
//			expectedResponse: delivery.NewResponse(delivery.StatusResponseSuccessful, userdelivery.ResponseSuccessfulLogOut),
//		},
//	}
//
//	for _, testCase := range testCases {
//		testCase := testCase
//
//		t.Run(testCase.name, func(t *testing.T) {
//			t.Parallel()
//
//			req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
//			req.Header.Add("Cookie", testCase.inputCookie.String())
//
//			w := httptest.NewRecorder()
//
//			authStorageMap := repository.NewAuthStorageMap()
//			authHandler := &userdelivery.AuthHandler{
//				storage: authStorageMap,
//			}
//
//			authHandler.LogOutHandler(w, req)
//
//			resp := w.Result()
//			defer resp.Body.Close()
//
//			receivedResponse, err := io.ReadAll(resp.Body)
//			if err != nil {
//				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
//			}
//
//			var resultResponse delivery.Response
//
//			err = json.Unmarshal(receivedResponse, &resultResponse)
//			if err != nil {
//				t.Fatalf("Failed to Unmarshal(receivedResponse): %v\n receivedResponse: %+v", err, receivedResponse)
//			}
//
//			if !reflect.DeepEqual(*testCase.expectedResponse, resultResponse) {
//				t.Errorf("wrong Response: got %+v, expected %+v",
//					resultResponse, testCase.expectedResponse)
//			}
//
//			allCookies := resp.Cookies()
//			for _, cookie := range allCookies {
//				if cookie.Name == userdelivery.CookieAuthName {
//					if cookie.Expires.Before(time.Now()) {
//						return
//					}
//				}
//			}
//
//			t.Fatalf("wrong cookie expire: %+v", allCookies)
//		})
//	}
//}
