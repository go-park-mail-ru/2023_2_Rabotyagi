package usecases_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
)

func NewCommentService(ctrl *gomock.Controller,
	behaviorCommentStorage func(m *mocks.MockICommentStorage),
) (*usecases.CommentService, error) {
	_ = mylogger.NewNop()

	mockCommentService := mocks.NewMockICommentStorage(ctrl)

	behaviorCommentStorage(mockCommentService)

	commentService, err := usecases.NewCommentService(mockCommentService)
	if err != nil {
		return nil, fmt.Errorf("unexpected err=%w", err)
	}

	return commentService, nil
}

func TestAddComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()

	type TestCase struct {
		name                   string
		inputReader            io.Reader
		behaviorCommentStorage func(m *mocks.MockICommentStorage)
		expectedCommentID      uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"sender_id":1,
					"recipient_id": 2, 
					"rating": 4,
					"text": "good"}`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().AddComment(baseCtx, &models.PreComment{
					SenderID:    uint64(1),
					RecipientID: uint64(2), Rating: uint8(4), Text: "good",
				}).Return(
					uint64(1), nil)
			},
			expectedCommentID: uint64(1),
			expectedError:     nil,
		},
		{
			name: "test validation error",
			inputReader: strings.NewReader(
				`{"rating": 4}`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {},
			expectedCommentID:      uint64(0),
			expectedError:          usecases.ErrValidatePreComment,
		},
		{
			name: "test validation error",
			inputReader: strings.NewReader(
				`{"sender_id":1,
					"recipient_id": 1, 
					"rating": 4,
					"text": "good"}`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {},
			expectedCommentID:      uint64(0),
			expectedError:          usecases.ErrCommentingYourself,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commentService, err := NewCommentService(ctrl, testCase.behaviorCommentStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			commentID, err := commentService.AddComment(baseCtx, testCase.inputReader, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(commentID, testCase.expectedCommentID); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type testCase struct {
		name                   string
		inputCommentID         uint64
		behaviorCommentStorage func(m *mocks.MockICommentStorage)
		expectedError          error
	}

	testCases := [...]testCase{
		{
			name:           "test basic work",
			inputCommentID: 1,
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().DeleteComment(baseCtx, uint64(1), test.UserID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:           "test internal error",
			inputCommentID: 1,
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().DeleteComment(baseCtx, uint64(1), test.UserID).Return(testInternalErr)
			},
			expectedError: testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewCommentService(ctrl, testCase.behaviorCommentStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.DeleteComment(baseCtx, testCase.inputCommentID, test.UserID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		inputReader            io.Reader
		behaviorCommentStorage func(m *mocks.MockICommentStorage)
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(
				`{"rating": 2, 
					"text": "bad" }`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().UpdateComment(baseCtx, test.UserID, test.CommentID,
					utils.StructToMap(models.CommentChanges{Text: "bad", Rating: 2})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "test update only rating",
			inputReader: strings.NewReader(
				`{"rating": 2}`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().UpdateComment(baseCtx, test.UserID, test.CommentID,
					utils.StructToMap(models.CommentChanges{Rating: 2})).Return(nil) //nolint:exhaustruct
			},
			expectedError: nil,
		},
		{
			name: "test update only text",
			inputReader: strings.NewReader(
				`{"text": "bad" }`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().UpdateComment(baseCtx, test.UserID, test.CommentID,
					utils.StructToMap(models.CommentChanges{Text: "bad"})).Return(nil) //nolint:exhaustruct
			},
			expectedError: nil,
		},
		{
			name: "test internal error",
			inputReader: strings.NewReader(
				`{"rating": 2, 
					"text": "bad" }`),
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().UpdateComment(baseCtx, test.UserID, test.CommentID,
					utils.StructToMap(models.CommentChanges{Text: "bad", Rating: 2})).Return(testInternalErr)
			},
			expectedError: testInternalErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewCommentService(ctrl, testCase.behaviorCommentStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			err = productService.UpdateComment(baseCtx, testCase.inputReader, test.UserID, test.CommentID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}
		})
	}
}

func TestGetCommentList(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	baseCtx := context.Background()
	testInternalErr := myerrors.NewErrorInternal("Test error")

	type TestCase struct {
		name                   string
		behaviorCommentStorage func(m *mocks.MockICommentStorage)
		expectedCommentInFeed  []*models.CommentInFeed
		expectedError          error
		senderID               uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().GetCommentList(baseCtx, uint64(1), uint64(2), test.UserID, uint64(2)).Return(
					[]*models.CommentInFeed{
						{
							ID: test.CommentID, SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
							Text: "good", Rating: 5, CreatedAt: time.Time{},
						},
					}, nil)
			},
			expectedCommentInFeed: []*models.CommentInFeed{
				{
					ID: test.CommentID, SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
					Text: "good", Rating: 5, CreatedAt: time.Time{},
				},
			},
			expectedError: nil,
			senderID:      uint64(2),
		},
		{
			name: "test internal error",
			behaviorCommentStorage: func(m *mocks.MockICommentStorage) {
				m.EXPECT().GetCommentList(baseCtx, uint64(1), uint64(2), test.UserID, uint64(2)).Return(
					nil, testInternalErr)
			},
			expectedCommentInFeed: nil,
			expectedError:         testInternalErr,
			senderID:              uint64(2),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			productService, err := NewCommentService(ctrl, testCase.behaviorCommentStorage)
			if err != nil {
				t.Fatalf("Failed create productService %+v", err)
			}

			ordersInBasket, err := productService.GetCommentList(baseCtx, uint64(1), test.CountComment,
				test.UserID, testCase.senderID)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(ordersInBasket, testCase.expectedCommentInFeed); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}
