package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAddComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		preComment             *models.PreComment
		expectedResponse       uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`INSERT INTO public."comment"`).WithArgs(uint64(2),
					uint64(1), "good", uint8(5)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectQuery(`SELECT last_value FROM "public"."comment_id_seq";`).
					WillReturnRows(pgxmock.NewRows([]string{"last_value"}).
						AddRow(uint64(1)))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			preComment:       &models.PreComment{SenderID: uint64(1), RecipientID: uint64(2), Rating: uint8(5), Text: "good"},
			expectedResponse: uint64(1),
		},
		{
			name: "test 0 rows affected",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`INSERT INTO public."comment"`).WithArgs(uint64(2),
					uint64(1), "", uint8(5)).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))

				mockPool.ExpectQuery(`SELECT last_value FROM "public"."comment_id_seq";`).
					WillReturnRows(pgxmock.NewRows([]string{"last_value"}).
						AddRow(uint64(1)))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			preComment:       &models.PreComment{SenderID: uint64(1), RecipientID: uint64(2), Rating: 5, Text: ""},
			expectedResponse: uint64(1),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage, mockPool)

			response, err := catStorage.AddComment(ctx, testCase.preComment)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		userID                 uint64
		commentID              uint64
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."comment"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:        1,
			commentID:     1,
			expectedError: nil,
		},
		{
			name: "test no affected rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."comment"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			userID:        1,
			commentID:     1,
			expectedError: repository.ErrNoAffectedCommentRows,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage, mockPool)

			errActual := commentStorage.DeleteComment(ctx, testCase.commentID, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualError(errActual, testCase.expectedError)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetCommentList(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		resID                  uint64
		senderID               uint64
		offset                 uint64
		count                  uint64
		expectedResponse       []*models.CommentInFeed
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT c.id AS comment_id, c.sender_id, 
       CASE WHEN u.name IS NOT NULL THEN u.name ELSE u.email END,
       u.avatar,
       c.text,
       c.rating,
       c.created_at
FROM public."comment" c`).WithArgs(uint64(1), uint64(2), uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"comment_id", "sender_id", "name", "avatar",
						"text", "rating", "created_at",
					}).
						AddRow(uint64(1), uint64(2), "Ivan", sql.NullString{Valid: false, String: ""}, "good", uint8(5), time.Time{}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			resID:    1,
			senderID: 2,
			offset:   1,
			count:    1,
			expectedResponse: []*models.CommentInFeed{
				{
					ID: 1, SenderID: 2, SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
					Text: "good", Rating: 5, CreatedAt: time.Time{},
				},
			},
		},
		{
			name: "test empty",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT c.id AS comment_id, c.sender_id, 
       CASE WHEN u.name IS NOT NULL THEN u.name ELSE u.email END,
       u.avatar,
       c.text,
       c.rating,
       c.created_at
FROM public."comment" c`).WithArgs(uint64(1), uint64(2), uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			resID:            1,
			senderID:         2,
			offset:           1,
			count:            1,
			expectedResponse: nil,
		},
		{
			name: "test more rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectQuery(`SELECT c.id AS comment_id, c.sender_id, 
       CASE WHEN u.name IS NOT NULL THEN u.name ELSE u.email END,
       u.avatar,
       c.text,
       c.rating,
       c.created_at
FROM public."comment" c`).WithArgs(uint64(1), uint64(2), uint64(1), uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{
						"comment_id", "sender_id", "name", "avatar",
						"text", "rating", "created_at",
					}).
						AddRow(uint64(1), uint64(2), "Ivan", sql.NullString{Valid: false, String: ""}, "good", uint8(5), time.Time{}).
						AddRow(uint64(2), uint64(3), "Petr", sql.NullString{Valid: false, String: ""}, "bad", uint8(2), time.Time{}).
						AddRow(uint64(3), uint64(4), "Mark", sql.NullString{Valid: false, String: ""}, "not bad", uint8(3), time.Time{}))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			resID:    1,
			senderID: 2,
			offset:   1,
			count:    1,
			expectedResponse: []*models.CommentInFeed{
				{
					ID: 1, SenderID: 2, SenderName: "Ivan", Avatar: sql.NullString{Valid: false, String: ""},
					Text: "good", Rating: 5, CreatedAt: time.Time{},
				},
				{
					ID: 2, SenderID: 3, SenderName: "Petr", Avatar: sql.NullString{Valid: false, String: ""},
					Text: "bad", Rating: 2, CreatedAt: time.Time{},
				},
				{
					ID: 3, SenderID: 4, SenderName: "Mark", Avatar: sql.NullString{Valid: false, String: ""},
					Text: "not bad", Rating: 3, CreatedAt: time.Time{},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage, mockPool)

			response, err := commentStorage.GetCommentList(ctx, testCase.offset, testCase.count,
				testCase.resID, testCase.senderID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualTest(response, testCase.expectedResponse)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCommentUpdate(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface)
		userID                 uint64
		commentID              uint64
		updateFields           map[string]interface{}
		expectedError          error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."comment"`).WithArgs(uint8(3), "not bad", uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:        1,
			commentID:     1,
			updateFields:  map[string]interface{}{"rating": uint8(3), "text": "not bad"},
			expectedError: nil,
		},
		{
			name: "test no rows",
			behaviorProductStorage: func(m *repository.ProductStorage, mockPool pgxmock.PgxPoolIface) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."comment"`).WithArgs(uint8(3), "not bad", uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))

				mockPool.ExpectRollback()
				mockPool.ExpectRollback()
			},
			userID:        1,
			commentID:     1,
			updateFields:  map[string]interface{}{"rating": uint8(3), "text": "not bad"},
			expectedError: repository.ErrNoAffectedCommentRows,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockPool, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("%v", err)
			}

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage, mockPool)

			errActual := commentStorage.UpdateComment(ctx, testCase.userID, testCase.commentID, testCase.updateFields)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			err = utils.EqualError(errActual, testCase.expectedError)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
