package delivery_test

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/mocks"
	mocksauth "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/auth/mocks"

	"go.uber.org/mock/gomock"
)

func NewProfileHandler(ctrl *gomock.Controller,
	behaviorUserService func(m *mocks.MockIUserService),
	behaviorSessionManagerClient func(m *mocksauth.MockSessionMangerClient),
) (*delivery.ProfileHandler, error) {
	mockUserService := mocks.NewMockIUserService(ctrl)
	mockSessionManagerClient := mocksauth.NewMockSessionMangerClient(ctrl)

	behaviorUserService(mockUserService)
	behaviorSessionManagerClient(mockSessionManagerClient)

	profileHandler, err := delivery.NewProfileHandler(mockUserService, mockSessionManagerClient)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return profileHandler, nil
}

//nolint:funlen
func TestPartiallyUpdateUser(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                         string
		behaviorUserService          func(m *mocks.MockIUserService)
		behaviorSessionManagerClient func(m *mocksauth.MockSessionMangerClient)
		request                      *http.Request
		expectedResponse             any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work patch",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, "/api/v1/profile/update", strings.NewReader(
					`{"email":"new_email@mail.ru"}`))
				req.AddCookie(&test.Cookie) //nolint:exhaustruct

				return req
			}(),
			behaviorUserService: func(m *mocks.MockIUserService) {
				m.EXPECT().UpdateUser(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"email":"new_email@mail.ru"}`)), true, test.UserID).Return(
					&models.UserWithoutPassword{ //nolint:exhaustruct
						ID:    test.UserID,
						Email: "new_email@mail.ru", Name: sql.NullString{Valid: true, String: "test_name"},
					}, nil)
			},
			behaviorSessionManagerClient: func(m *mocksauth.MockSessionMangerClient) {
				m.EXPECT().Check(gomock.Any(), &auth.Session{AccessToken: test.AccessToken}).Return(
					&auth.UserID{UserId: test.UserID}, nil)
			},
			expectedResponse: delivery.NewProfileResponse(&models.UserWithoutPassword{ //nolint:exhaustruct
				ID:    test.UserID,
				Email: "new_email@mail.ru", Name: sql.NullString{Valid: true, String: "test_name"},
			}),
		},
		{
			name:                         "test cookie not presented",
			request:                      httptest.NewRequest(http.MethodPatch, "/api/v1/profile/update", nil),
			behaviorUserService:          func(m *mocks.MockIUserService) {},
			behaviorSessionManagerClient: func(m *mocksauth.MockSessionMangerClient) {},
			expectedResponse: responses.NewErrResponse(
				responses.ErrCookieNotPresented.Status(), responses.ErrCookieNotPresented.Error()),
		},
		{
			name: "test internal in SessionManagerClient",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, "/api/v1/profile/update", nil)
				req.AddCookie(&test.Cookie) //nolint:exhaustruct

				return req
			}(),
			behaviorUserService: func(m *mocks.MockIUserService) {},
			behaviorSessionManagerClient: func(m *mocksauth.MockSessionMangerClient) {
				m.EXPECT().Check(gomock.Any(), &auth.Session{AccessToken: test.AccessToken}).Return(
					nil, myerrors.NewErrorInternal("Test error"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:                         "test method not allowed",
			request:                      httptest.NewRequest(http.MethodDelete, "/api/v1/profile/update", nil),
			behaviorUserService:          func(m *mocks.MockIUserService) {},
			behaviorSessionManagerClient: func(m *mocksauth.MockSessionMangerClient) {},
			expectedResponse: `Method not allowed
`,
		},
		{
			name: "test error put internal error",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/profile/update", strings.NewReader(
					`{"email":"new_email@mail.ru, "name":"test"}`))
				req.AddCookie(&test.Cookie) //nolint:exhaustruct

				return req
			}(),
			behaviorUserService: func(m *mocks.MockIUserService) {
				m.EXPECT().UpdateUser(gomock.Any(), io.NopCloser(strings.NewReader(
					`{"email":"new_email@mail.ru, "name":"test"}`)), false, test.UserID).Return(
					nil, myerrors.NewErrorInternal("Test error"))
			},
			behaviorSessionManagerClient: func(m *mocksauth.MockSessionMangerClient) {
				m.EXPECT().Check(gomock.Any(), &auth.Session{AccessToken: test.AccessToken}).Return(
					&auth.UserID{UserId: test.UserID}, nil)
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileHandler, err := NewProfileHandler(ctrl, testCase.behaviorUserService, testCase.behaviorSessionManagerClient)
			if err != nil {
				t.Fatalf("Failed create profileHandler %s", err.Error())
			}

			w := httptest.NewRecorder()

			profileHandler.PartiallyUpdateUserHandler(w, testCase.request)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}
