package delivery

import (
	"database/sql"
	"errors"
	mock_category "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/mocks"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDelivery_ListByUser(t *testing.T) {
	type fields struct {
		serv *mock_category.MockICategoryService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   http.ServeMux
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetFullCategories([]*models.Category{
					{ID: 1, Name: "aaaa", ParentID: sql.NullInt64{Int64: 0, Valid: false}},
					{ID: 2, Name: "bbbb", ParentID: sql.NullInt64{Int64: 1, Valid: true}},
					{ID: 3, Name: "aaaa", ParentID: sql.NullInt64{Valid: true, Int64: 2}},
				})
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "2"}},
			response: `{"items":[{"id":2,"user1_id":2,"user2_id":3,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":3,"user1_id":8,"user2_id":2,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":4,"user1_id":2,"user2_id":4,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`,
			err:      nil,
		},
		"no chats": {
			prepare: func(f *fields) {
				f.serv.EXPECT().ListByUser(3).Return([]models.Chat{}, nil)
			},
			params:   []httprouter.Param{{Key: "user-id", Value: "3"}},
			response: `{"items":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "user-id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{
				serv: f.serv,
				log:  logger,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			err = del.ListByUser(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected: %s\nGot: %s", test.response, string(body))
			}
		})
	}
}
