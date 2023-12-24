package delivery_test

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

func TestAddComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/comment/add", strings.NewReader(
				`{
    "recipient_id": 1,
    "text": "good",
    "rating": 4
}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddComment(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
    "recipient_id": 1,
    "text": "good",
    "rating": 4
}`)), test.UserID).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
		{
			name: "test another data",
			request: httptest.NewRequest(http.MethodPost, "/api/v1/comment/add", strings.NewReader(
				`{
    "recipient_id": 1,
    "text": "bad",
    "rating": 1
}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().AddComment(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
    "recipient_id": 1,
    "text": "bad",
    "rating": 1
}`)), test.UserID).Return(uint64(1), nil)
			},
			expectedResponse: responses.ResponseID{
				Status: statuses.StatusRedirectAfterSuccessful,
				Body:   responses.ResponseBodyID{ID: 1},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&test.Cookie)
			productHandler.AddCommentHandler(w, testCase.request)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryID                string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().DeleteComment(gomock.Any(), uint64(1), test.UserID).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulDeleteComment},
			},
		},
		{
			name:    "test error in internal",
			queryID: "1",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().DeleteComment(gomock.Any(), uint64(1), test.UserID).Return(
					myerrors.NewErrorInternal("Test Error Internal"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name:    "test error uncorrected query param",
			queryID: "wrong type",
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest,
				fmt.Sprintf("%s comment_id=wrong type", utils.MessageErrWrongNumberParam)),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/comment/delete", nil)
			utils.AddQueryParamsToRequest(req, map[string]string{"comment_id": testCase.queryID})
			req.AddCookie(&test.Cookie)
			productHandler.DeleteCommentHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetCommentList(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		queryParams            map[string]string
		behaviorProductService func(m *mocks.MockIProductService)
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			queryParams: map[string]string{"count": "2", "offset": "1", "user_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetCommentList(gomock.Any(), uint64(1), uint64(2), test.UserID, uint64(1)).Return(
					[]*models.CommentInFeed{
						{
							ID: 1, SenderID: uint64(2), SenderName: "Ivan",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
						{
							ID: 2, SenderID: uint64(3), SenderName: "Petr",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
					}, nil)
			},
			expectedResponse: delivery.NewCommentListResponse(
				[]*models.CommentInFeed{
					{
						ID: 1, SenderID: uint64(2), SenderName: "Ivan",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
					{
						ID: 2, SenderID: uint64(3), SenderName: "Petr",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
				}),
		},
		{
			name:        "test zero work",
			queryParams: map[string]string{"count": "0", "offset": "0", "user_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetCommentList(gomock.Any(), uint64(0), uint64(0), test.UserID, uint64(1)).Return(
					[]*models.CommentInFeed{}, nil)
			},
			expectedResponse: delivery.NewCommentListResponse(
				[]*models.CommentInFeed{}),
		},
		{
			name:        "test a lot of count",
			queryParams: map[string]string{"count": "5", "offset": "1", "user_id": "1"},
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().GetCommentList(gomock.Any(), uint64(1), uint64(5), test.UserID, uint64(1)).Return(
					[]*models.CommentInFeed{ //nolint:dupl
						{
							ID: 1, SenderID: uint64(2), SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
							Text: "Good", Rating: uint8(5), CreatedAt: time.Time{},
						},
						{
							ID: 2, SenderID: uint64(3), SenderName: "Petr",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
						{
							ID: 3, SenderID: uint64(4), SenderName: "Petr",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
						{
							ID: 4, SenderID: uint64(5), SenderName: "Petr",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
						{
							ID: 5, SenderID: uint64(6), SenderName: "Petr",
							Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
							CreatedAt: time.Time{},
						},
					}, nil)
			},
			expectedResponse: delivery.NewCommentListResponse(
				[]*models.CommentInFeed{ //nolint:dupl
					{
						ID: 1, SenderID: uint64(2), SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
						Text: "Good", Rating: uint8(5), CreatedAt: time.Time{},
					},
					{
						ID: 2, SenderID: uint64(3), SenderName: "Petr",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
					{
						ID: 3, SenderID: uint64(4), SenderName: "Petr",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
					{
						ID: 4, SenderID: uint64(5), SenderName: "Petr",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
					{
						ID: 5, SenderID: uint64(6), SenderName: "Petr",
						Avatar: sql.NullString{Valid: false, String: ""}, Text: "Good", Rating: uint8(5),
						CreatedAt: time.Time{},
					},
				}),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/comment/get_list", nil)
			utils.AddQueryParamsToRequest(req, testCase.queryParams)
			req.AddCookie(&test.Cookie)
			productHandler.GetCommentListHandler(recorder, req)

			err = test.CompareHTTPTestResult(recorder, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockIProductService)
		request                *http.Request
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/comment/update?comment_id=1", strings.NewReader(
				`{
    "rating":3,
    "text":"not bad"
}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateComment(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
    "rating":3,
    "text":"not bad"
}`)), test.UserID, uint64(1)).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulUpdateComment},
			},
		},
		{
			name: "test update only rating",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/comment/update?comment_id=1", strings.NewReader(
				`{
    "rating":3
}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateComment(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
    "rating":3
}`)), test.UserID, uint64(1)).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulUpdateComment},
			},
		},
		{
			name: "test update only text",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/comment/update?comment_id=1", strings.NewReader(
				`{
    "text":"not bad"
}`)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateComment(gomock.Any(), io.NopCloser(strings.NewReader(
					`{
    "text":"not bad"
}`)), test.UserID, uint64(1)).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulUpdateComment},
			},
		},
		{
			name: "test empty patch",
			request: httptest.NewRequest(http.MethodPatch, "/api/v1/comment/update?comment_id=1", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT().UpdateComment(gomock.Any(), io.NopCloser(strings.NewReader(
					``)), test.UserID, uint64(1)).Return(nil)
			},
			expectedResponse: responses.ResponseSuccessful{
				Status: statuses.StatusResponseSuccessful,
				Body:   responses.ResponseBody{Message: delivery.ResponseSuccessfulUpdateComment},
			},
		},

		{
			name: "test wrong method",
			request: httptest.NewRequest(http.MethodDelete, "/api/v1/comment/update", strings.NewReader(
				``)),
			behaviorProductService: func(m *mocks.MockIProductService) {
				m.EXPECT()
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

			productHandler, err := NewProductHandler(ctrl, testCase.behaviorProductService)
			if err != nil {
				t.Fatalf("Failed create productHandler %+v", err)
			}

			w := httptest.NewRecorder()

			testCase.request.AddCookie(&test.Cookie)
			productHandler.UpdateCommentHandler(w, testCase.request)

			err = test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}
