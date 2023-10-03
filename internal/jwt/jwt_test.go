package jwt_test

import (
	"errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/jwt"
	"reflect"
	"testing"
)

//nolint:funlen
func TestGenerateJwtToken(t *testing.T) {
	t.Parallel()

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
			inputUserJwtPayload: &jwt.UserJwtPayload{UserID: 1, Expire: expire, Email: "example@mail.ru"},
			expectedJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJlbWFpbCI6ImV4YW1wbGVAbWFpbC5ydSIsImV4cGlyZSI6MCwidXNlcklEIjoxfQ." +
				"GBCEb3XJ6aHTsyl8jC3lxSWK6byjbYN0kg2e3NH2i9s",
			expectedError: nil,
		},
		{
			name:                "test empty UserJwtPayload",
			inputUserJwtPayload: &jwt.UserJwtPayload{UserID: 0, Expire: expire, Email: ""},
			expectedJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImV4cGlyZSI6MCwidXNlcklEIjowfQ." +
				"s4EqX-V9Q3pWejqJe0x8Z65PZFVtzeu3ByV8txPboTo",
			expectedError: nil,
		},
		{
			name: "test long email and big userID",
			inputUserJwtPayload: &jwt.UserJwtPayload{
				UserID: 100000,
				Expire: expire,
				Email:  "exampleaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@mail.ru",
			},
			expectedJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJlbWFpbCI6ImV4YW1wbGVhYWFhYWFhYWFhYWFhYWF" +
				"hYWFhYWFhYWFhYWFhYWFhYWFhYWFhQG1haWwucnUiLCJleHBpcmUiOjAsInVzZXJJRCI6MTAwMDAwfQ." +
				"CQIXkeDEW3Y0ffLm9efgsozkWvLK1sg4ArmYBReHjsE",
			expectedError: nil,
		},
		{
			name:                "test nil UserJwtPayload",
			inputUserJwtPayload: nil,
			expectedJwt:         "",
			expectedError:       jwt.ErrNilToken,
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
			expectedUserJwtPayload: &jwt.UserJwtPayload{UserID: 1, Expire: 0, Email: "example@mail.ru"},
			expectedError:          nil,
		},
		{
			name: "test empty UserJwtPayload",
			inputRawJwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImV4cGlyZSI6MCwidXNlcklEIjowfQ." +
				"s4EqX-V9Q3pWejqJe0x8Z65PZFVtzeu3ByV8txPboTo",
			expectedUserJwtPayload: &jwt.UserJwtPayload{UserID: 0, Expire: 0, Email: ""},
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
				Email:  "exampleaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@mail.ru",
			},
			expectedError: nil,
		},
		{
			name:                   "test parsing invalid token",
			inputRawJwt:            "",
			expectedUserJwtPayload: nil,
			expectedError:          jwt.ErrParseToken,
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
