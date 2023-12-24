package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/product/repository"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/pashagolub/pgxmock/v3"
)

func TestAddComment(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		preComment             *models.PreComment
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
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
			preComment: &models.PreComment{SenderID: uint64(1), RecipientID: uint64(2), Rating: uint8(5), Text: "good"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			catStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(catStorage)

			_, err = catStorage.AddComment(ctx, testCase.preComment)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) { //nolint:dupl
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
		commentID              uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`DELETE FROM public."comment"`).WithArgs(uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:    1,
			commentID: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage)

			err = commentStorage.DeleteComment(ctx, testCase.commentID, testCase.userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetCommentList(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		resID                  uint64
		senderID               uint64
		offset                 uint64
		count                  uint64
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
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
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage)

			_, err = commentStorage.GetCommentList(ctx, testCase.offset, testCase.count, testCase.resID, testCase.senderID)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCommentUpdate(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("%v", err)
	}

	type TestCase struct {
		name                   string
		behaviorProductStorage func(m *repository.ProductStorage)
		userID                 uint64
		commentID              uint64
		updateFields           map[string]interface{}
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			behaviorProductStorage: func(m *repository.ProductStorage) {
				mockPool.ExpectBegin()

				mockPool.ExpectExec(`UPDATE public."comment"`).WithArgs(uint8(3), "not bad", uint64(1), uint64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()
			},
			userID:       1,
			commentID:    1,
			updateFields: map[string]interface{}{"rating": uint8(3), "text": "not bad"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			commentStorage, err := repository.NewProductStorage(mockPool)
			if err != nil {
				t.Fatalf("%v", err)
			}

			testCase.behaviorProductStorage(commentStorage)

			err = commentStorage.UpdateComment(ctx, testCase.userID, testCase.commentID, testCase.updateFields)
			if err != nil {
				t.Fatal(err)
			}

			if err := mockPool.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
