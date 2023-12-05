package usecases_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/user/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

//nolint:nolintlint,funlen
func TestValidateUserWithoutPassword(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                        string
		inputReader                 io.Reader
		expectedUserWithoutPassword *models.UserWithoutPassword
		expectedError               error
	}

	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputReader: strings.NewReader(`{"id":1, "email":"ivn-tyt@mail.ru",
"created_at":"2000-01-01T00:00:00Z"}`),
			expectedUserWithoutPassword: &models.UserWithoutPassword{ //nolint:exhaustruct
				ID: 1, Email: "ivn-tyt@mail.ru",
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: nil,
		},
		{
			name: "test error validation",
			inputReader: strings.NewReader(`{"id":1, "email":"not_email",
"created_at":"2000-01-01T00:00:00Z"}`),
			expectedUserWithoutPassword: nil,
			expectedError:               usecases.ErrValidateUserWithoutPassword,
		},
		{
			name:                        "test error decode",
			inputReader:                 strings.NewReader(`{`),
			expectedUserWithoutPassword: nil,
			expectedError:               usecases.ErrDecodeUser,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			userWithoutPassword, err := usecases.ValidateUserWithoutPassword(testCase.inputReader)
			if errInner := utils.EqualError(err, testCase.expectedError); errInner != nil {
				t.Fatalf("Failed EqualError: %+v", errInner)
			}

			if err := utils.EqualTest(userWithoutPassword, testCase.expectedUserWithoutPassword); err != nil {
				t.Fatalf("Failed EqualTest %+v", err)
			}
		})
	}
}
