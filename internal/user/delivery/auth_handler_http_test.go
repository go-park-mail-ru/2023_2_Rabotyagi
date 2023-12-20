package delivery_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

func NewAuthHandler(ctrl *gomock.Controller,
	behaviorSessionManagerClient func(m *mocks.MockSessionMangerClient),
) (*delivery.AuthHandler, error) {
	mockSessionManagerClient := mocks.NewMockSessionMangerClient(ctrl)

	behaviorSessionManagerClient(mockSessionManagerClient)

	authHandler, err := delivery.NewAuthHandler(mockSessionManagerClient)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return authHandler, nil
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                         string
		behaviorSessionManagerClient func(m *mocks.MockSessionMangerClient)
		request                      *http.Request
		checkHeader                  func(recorder *httptest.ResponseRecorder) error
		expectedResponse             any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(
				`{"email":"ivn-tyt@mail.ru", "password": "strong"}`)),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Create(gomock.Any(),
					&auth.User{Email: "ivn-tyt@mail.ru", Password: "strong"}).Return(
					&auth.Session{AccessToken: test.AccessToken}, nil)
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				cookieRaw := recorder.Header().Get("Set-Cookie")
				if !strings.Contains(cookieRaw, test.AccessToken) {
					return fmt.Errorf("cookie not contain jwt token. Cookie: %s", cookieRaw) //nolint
				}

				return nil
			},
			expectedResponse: responses.NewResponseSuccessful(delivery.ResponseSuccessfulSignUp),
		},
		{
			name: "test internal error",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(
				`{"email":"ivn-tyt@mail.ru", "password": "strong"}`)),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Create(gomock.Any(),
					&auth.User{Email: "ivn-tyt@mail.ru", Password: "strong"}).Return(
					nil, myerrors.NewErrorInternal("Test error"))
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name: "test wrong method",
			request: httptest.NewRequest(http.MethodDelete, "/api/v1/signup", strings.NewReader(
				`{"email":"ivn-tyt@mail.ru", "password": "strong"}`)),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: `Method not allowed
`,
		},
		{
			name: "test wrong format",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(
				`{`)),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest, "unexpected EOF"),
		},
		{
			name: "test wrong content",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(
				`{"email":"wrong_email", "password":"123333"}`)),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadContentRequest, "Некорректный формат email"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authHandler, err := NewAuthHandler(ctrl, testCase.behaviorSessionManagerClient)
			if err != nil {
				t.Fatalf("Failed create authHandler %s", err.Error())
			}

			recorder := httptest.NewRecorder()

			authHandler.SignUpHandler(recorder, testCase.request)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}

			err = testCase.checkHeader(recorder)
			if err != nil {
				t.Fatalf("Wrong Headers %s", err.Error())
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                         string
		behaviorSessionManagerClient func(m *mocks.MockSessionMangerClient)
		request                      *http.Request
		checkHeader                  func(recorder *httptest.ResponseRecorder) error
		expectedResponse             any
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/signin?email=ivn-tyt@mail.ru&password=strong", nil),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Login(gomock.Any(),
					&auth.User{Email: "ivn-tyt@mail.ru", Password: "strong"}).Return(
					&auth.Session{AccessToken: test.AccessToken}, nil)
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				cookieRaw := recorder.Header().Get("Set-Cookie")
				if !strings.Contains(cookieRaw, test.AccessToken) {
					return fmt.Errorf("cookie not contain jwt token. Cookie: %s", cookieRaw) //nolint
				}

				return nil
			},
			expectedResponse: responses.NewResponseSuccessful(delivery.ResponseSuccessfulSignIn),
		},
		{
			name: "test internal error",
			request: httptest.NewRequest(http.MethodGet,
				"/api/v1/signin?email=ivn-tyt@mail.ru&password=strong", nil),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Login(gomock.Any(),
					&auth.User{Email: "ivn-tyt@mail.ru", Password: "strong"}).Return(
					nil, myerrors.NewErrorInternal("Test error"))
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name: "test wrong method",
			request: httptest.NewRequest(http.MethodDelete,
				"/api/v1/signin", nil),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: `Method not allowed
`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authHandler, err := NewAuthHandler(ctrl, testCase.behaviorSessionManagerClient)
			if err != nil {
				t.Fatalf("Failed create authHandler %s", err.Error())
			}

			recorder := httptest.NewRecorder()

			authHandler.SignInHandler(recorder, testCase.request)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}

			err = testCase.checkHeader(recorder)
			if err != nil {
				t.Fatalf("Wrong Headers %s", err.Error())
			}
		})
	}
}

func TestLogOut(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                         string
		behaviorSessionManagerClient func(m *mocks.MockSessionMangerClient)
		request                      *http.Request
		checkHeader                  func(recorder *httptest.ResponseRecorder) error
		expectedResponse             any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
				req.AddCookie(&test.Cookie)

				return req
			}(),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Delete(gomock.Any(),
					&auth.Session{AccessToken: test.AccessToken}).Return(
					&auth.Session{AccessToken: "jwt_test_token"}, nil)
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				cookieRaw := recorder.Header().Get("Set-Cookie")
				if !strings.Contains(cookieRaw, "jwt_test_token") {
					return fmt.Errorf("cookie not contain jwt token. Cookie: %s", cookieRaw) //nolint
				}

				return nil
			},
			expectedResponse: responses.NewResponseSuccessful(delivery.ResponseSuccessfulLogOut),
		},
		{
			name: "test internal error",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
				req.AddCookie(&test.Cookie)

				return req
			}(),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {
				m.EXPECT().Delete(gomock.Any(),
					&auth.Session{AccessToken: test.AccessToken}).Return(nil,
					myerrors.NewErrorInternal("Test error"))
			},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:                         "test cookie not presented",
			request:                      httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: responses.NewErrResponse(
				responses.ErrCookieNotPresented.Status(), responses.ErrCookieNotPresented.Error()),
		},
		{
			name:                         "test wrong method",
			request:                      httptest.NewRequest(http.MethodDelete, "/api/v1/logout", nil),
			behaviorSessionManagerClient: func(m *mocks.MockSessionMangerClient) {},
			checkHeader: func(recorder *httptest.ResponseRecorder) error {
				return nil
			},
			expectedResponse: `Method not allowed
`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authHandler, err := NewAuthHandler(ctrl, testCase.behaviorSessionManagerClient)
			if err != nil {
				t.Fatalf("Failed create authHandler %s", err.Error())
			}

			recorder := httptest.NewRecorder()

			authHandler.LogOutHandler(recorder, testCase.request)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}

			err = testCase.checkHeader(recorder)
			if err != nil {
				t.Fatalf("Wrong Headers %s", err.Error())
			}
		})
	}
}
