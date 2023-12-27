package jwt_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/auth/internal/jwt"
)

func TestGenerateJwtToken(t *testing.T) {
	t.Parallel()

	mylogger.NewNop()

	type TestCase struct {
		name                string
		inputUserJwtPayload *jwt.UserJwtPayload
		expectedJwt         string
		expectedError       error
	}

	var expire int64

	secret := []byte("thisIsTestSecretItCanBeWeak")
	testCases := [...]TestCase{
		{
			name:                "test basic work",
			inputUserJwtPayload: &jwt.UserJwtPayload{UserID: 1, Expire: expire},
			expectedJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJleHBpcmUiOjAsInVzZXJJRCI6MX0." +
				"AxRUwsu1l_QGC2guNzeTYYOdIglczlCEyO2WUDKwfFw",
			expectedError: nil,
		},
		{
			name: "test big userID",
			inputUserJwtPayload: &jwt.UserJwtPayload{
				UserID: 100000,
				Expire: expire,
			},
			expectedJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJleHBpcmUiOjAsInVzZXJJRCI6MTAwMDAwfQ." +
				"TzJp8AoDA8hX1vQGc-m6OSAGO_w_n2gXu5J0ePJkK_k",
			expectedError: nil,
		},
		{
			name:                "test nil UserJwtPayload",
			inputUserJwtPayload: nil,
			expectedJwt:         "",
			expectedError:       jwt.ErrInvalidToken,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			receivedJwt, err := jwt.GenerateJwtToken(testCase.inputUserJwtPayload, secret)
			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("EXPECTED error: %v\n RECEIVED error: %v\n", testCase.expectedError, err)
			}

			if receivedJwt != testCase.expectedJwt {
				t.Errorf("EXPECTED jwt: %v\n RECEIVED jwt: %v\n", testCase.expectedJwt, receivedJwt)
			}
		})
	}
}

func TestNewUserJwtPayload(t *testing.T) {
	t.Parallel()

	mylogger.NewNop()

	type TestCase struct {
		name                   string
		inputRawJwt            string
		expectedUserJwtPayload *jwt.UserJwtPayload
		expectedError          error
	}

	secret := []byte("thisIsTestSecretItCanBeWeak")
	testCases := [...]TestCase{
		{
			name: "test basic work",
			inputRawJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJlbWFpbCI6ImV4YW1wbGVAbWFpbC5ydSIsImV4cGlyZSI6MCwidXNlcklEIjoxfQ." +
				"GBCEb3XJ6aHTsyl8jC3lxSWK6byjbYN0kg2e3NH2i9s",
			expectedUserJwtPayload: &jwt.UserJwtPayload{UserID: 1, Expire: 0},
			expectedError:          nil,
		},
		{
			name: "test long email and big userID",
			inputRawJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJlbWFpbCI6ImV4YW1wbGVhYWFhYWFhYWFhYWFhYWF" +
				"hYWFhYWFhYWFhYWFhYWFhYWFhYWFhQG1haWwucnUiLCJleHBpcmUiOjAsInVzZXJJRCI6MTAwMDAwfQ." +
				"CQIXkeDEW3Y0ffLm9efgsozkWvLK1sg4ArmYBReHjsE",
			expectedUserJwtPayload: &jwt.UserJwtPayload{
				UserID: 100000,
				Expire: 0,
			},
			expectedError: nil,
		},
		{
			name:                   "test parsing invalid token",
			inputRawJwt:            "",
			expectedUserJwtPayload: nil,
			expectedError:          jwt.ErrInvalidToken,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			receivedUserJwtPayload, err := jwt.NewUserJwtPayload(testCase.inputRawJwt, secret)
			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("EXPECTED error: %v\n RECEIVED error: %v\n", testCase.expectedError, err)
			}

			if !reflect.DeepEqual(receivedUserJwtPayload, testCase.expectedUserJwtPayload) {
				t.Errorf("EXPECTED UserJwtPayload: %v\n RECEIVED UserJwtPayload: %v\n",
					testCase.expectedUserJwtPayload, receivedUserJwtPayload)
			}
		})
	}
}

func TestSetSecret(t *testing.T) {
	t.Parallel()

	secret := []byte("test secret")

	jwt.SetSecret(secret)

	receivedSecret, err := jwt.GetSecret()
	if err != nil {
		t.Errorf("неожиданная ошибка %+v", err)
	}

	if err := utils.CompareSameType(string(receivedSecret), string(secret)); err != nil {
		t.Errorf("ошибка сравнения секретов %+v", err)
	}
}
