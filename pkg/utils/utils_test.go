package utils_test

import (
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils"
)

func TestSnakeToCamelCase(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		str      string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"Hello_World", "HelloWorld"},
		{"hello_world_test", "HelloWorldTest"},
		{"Hello_World_Test", "HelloWorldTest"},
		{"snake_case", "SnakeCase"},
		{"camel_case_test", "CamelCaseTest"},
	}

	for _, input := range inputs {
		actual := utils.SnakeToCamelCase(input.str)
		if actual != input.expected {
			t.Errorf("Input: %s, Expected: %s, Got: %s", input.str, input.expected, actual)
		}
	}
}

func TestCamelToSnakeCase(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		str      string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"helloWorldTest", "hello_world_test"},
		{"HelloWorldTest", "hello_world_test"},
		{"snake_case", "snake_case"},
		{"camelCaseTest", "camel_case_test"},
	}

	for _, input := range inputs {
		actual := utils.CamelToSnakeCase(input.str)
		if actual != input.expected {
			t.Errorf("Input: %s, Expected: %s, Got: %s", input.str, input.expected, actual)
		}
	}
}

func TestEqualTest(t *testing.T) {
	t.Parallel()

	type testCase struct {
		received any
		expected any
	}

	testCases := []testCase{
		{received: 5, expected: 5},
		{received: "str", expected: "str"},
		{received: "", expected: ""},
	}

	for _, tc := range testCases {
		err := utils.EqualTest(tc.received, tc.expected)
		if err != nil {
			t.Errorf("Received: %d, Expected: %d, Error: %s", tc.received, tc.expected, err.Error())
		}
	}
}

func TestHashPass(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name          string
		plainPassword string
	}

	testCases := [...]TestCase{
		{
			name:          "test basic work",
			plainPassword: "password",
		},
		{
			name:          "test empty",
			plainPassword: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			hashedPass, err := utils.HashPass(testCase.plainPassword)
			if err != nil {
				t.Errorf("Error hashing password: %s", err)
			}

			if hashedPass == "" {
				t.Errorf("Empty hashed password")
			}
		})
	}
}

func TestComparePassAndHash(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name          string
		plainPassword string
	}

	testCases := [...]TestCase{
		{
			name:          "test basic work",
			plainPassword: "password",
		},
		{
			name:          "test empty",
			plainPassword: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			hashedPass, err := utils.HashPass(testCase.plainPassword)
			if err != nil {
				t.Errorf("Error hashing password: %s", err)
			}

			hashedPassByteSl, err := hex.DecodeString(hashedPass)
			if err != nil {
				t.Errorf("Error hashing password: %s", err)
			}

			if !utils.ComparePassAndHash(hashedPassByteSl, testCase.plainPassword) {
				t.Errorf("Password and hash do not match")
			}

			if utils.ComparePassAndHash(hashedPassByteSl, "wrongpassword") {
				t.Errorf("Password and hash should not match")
			}
		})
	}
}

func TestHash256(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name         string
		pass         string
		expectedHash string
	}

	testCases := [...]TestCase{
		{
			name:         "basic",
			pass:         "48656c6c6f20476f7068657221",
			expectedHash: "be8c5fbcec1ca3f472cba2d613f780ae7c7efbaac657669adcf16a9cc525dd9b",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := hex.DecodeString(testCase.pass)
			if err != nil {
				t.Errorf("Error hashing content: %s", err)
			}

			hashedContent, err := utils.Hash256(decoded)
			if err != nil {
				t.Errorf("Error hashing content: %s", err)
			}

			if hashedContent != testCase.expectedHash {
				t.Errorf("Hash does not match. Expected: %s, Got: %s", testCase.expectedHash, hashedContent)
			}
		})
	}
}

func TestNullStringToUnsafe(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name        string
		validString string
	}

	testCases := [...]TestCase{
		{
			name:        "test basic work",
			validString: "password",
		},
		{
			name:        "test empty",
			validString: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			nullStr := sql.NullString{String: testCase.validString, Valid: true}

			result := utils.NullStringToUnsafe(nullStr)

			if result != nil {
				if *result != testCase.validString {
					t.Errorf("Expected %s, got %s", testCase.validString, *result)
				}
			} else {
				t.Errorf("Expected non-nil result, got nil")
			}
		})
	}
}

func TestUnsafeStringToNull(t *testing.T) {
	t.Parallel()

	nullStr := utils.UnsafeStringToNull(nil)

	if nullStr.Valid {
		t.Errorf("Expected Valid=false, got Valid=true")
	}

	validString := "valid string"
	nullStr = utils.UnsafeStringToNull(&validString)

	if !nullStr.Valid {
		t.Errorf("Expected Valid=true, got Valid=false")
	}

	if nullStr.String != validString {
		t.Errorf("Expected %s, got %s", validString, nullStr.String)
	}
}

func TestNullTimeToUnsafe(t *testing.T) {
	t.Parallel()

	validTime := time.Now()
	nullTime := sql.NullTime{Time: validTime, Valid: true}

	result := utils.NullTimeToUnsafe(nullTime)

	if result != nil {
		if *result != validTime {
			t.Errorf("Expected %s, got %s", validTime, *result)
		}
	} else {
		t.Errorf("Expected non-nil result, got nil")
	}
}

func TestUnsafeTimeToNull(t *testing.T) {
	t.Parallel()

	nullTime := utils.UnsafeTimeToNull(nil)

	if nullTime.Valid {
		t.Errorf("Expected Valid=false, got Valid=true")
	}

	validTime := time.Now()

	nullTime = utils.UnsafeTimeToNull(&validTime)
	if !nullTime.Valid {
		t.Errorf("Expected Valid=true, got Valid=false")
	}

	if nullTime.Time != validTime {
		t.Errorf("Expected %s, got %s", validTime, nullTime.Time)
	}
}

func TestNullInt64ToUnsafeUint(t *testing.T) {
	t.Parallel()

	validInt64 := int64(42)
	nullInt64 := sql.NullInt64{Int64: validInt64, Valid: true}

	result := utils.NullInt64ToUnsafeUint(nullInt64)

	if result != nil {
		if *result != uint64(validInt64) {
			t.Errorf("Expected %d, got %d", validInt64, *result)
		}
	} else {
		t.Errorf("Expected non-nil result, got nil")
	}
}

func TestUnsafeUint64ToNullInt(t *testing.T) {
	t.Parallel()

	nullInt64 := utils.UnsafeUint64ToNullInt(nil)
	if nullInt64.Valid {
		t.Errorf("Expected Valid=false, got Valid=true")
	}

	validUint64 := uint64(42)

	nullInt64 = utils.UnsafeUint64ToNullInt(&validUint64)
	if !nullInt64.Valid {
		t.Errorf("Expected Valid=true, got Valid=false")
	}

	if nullInt64.Int64 != int64(validUint64) {
		t.Errorf("Expected %d, got %d", validUint64, nullInt64.Int64)
	}
}

func TestAddQueryParamsToRequest(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name        string
		request     *http.Request
		params      map[string]string
		expectedURL string
	}

	testCases := [...]TestCase{
		{
			name:    "test basic work",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/product/get", nil),
			params: map[string]string{
				"product_id": "1",
			},
			expectedURL: "/api/v1/product/get?product_id=1",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			utils.AddQueryParamsToRequest(testCase.request, testCase.params)

			if testCase.request.URL.String() != testCase.expectedURL {
				t.Errorf("Expected URL: %s, got: %s", testCase.expectedURL, testCase.request.URL.String())
			}
		})
	}
}

func TestParseUint64FromRequest(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name          string
		request       *http.Request
		paramName     string
		expectedValue uint64
	}

	testCases := [...]TestCase{
		{
			name:          "test basic work",
			request:       httptest.NewRequest(http.MethodGet, "/api/v1/product/get?id=3", nil),
			paramName:     "id",
			expectedValue: 3,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			value, err := utils.ParseUint64FromRequest(testCase.request, testCase.paramName)
			if err != nil {
				t.Errorf("error while func calling %s", err)
			}

			if value != testCase.expectedValue {
				t.Errorf("Expected: %d, got: %d", testCase.expectedValue, value)
			}
		})
	}
}

func TestParseStringFromRequest(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	type TestCase struct {
		name          string
		request       *http.Request
		paramName     string
		expectedValue string
	}

	testCases := [...]TestCase{
		{
			name:          "test basic work",
			request:       httptest.NewRequest(http.MethodGet, "/api/v1/category/get?name=cars", nil),
			paramName:     "name",
			expectedValue: "cars",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			value := utils.ParseStringFromRequest(testCase.request, testCase.paramName)

			if value != testCase.expectedValue {
				t.Errorf("Expected: %s, got: %s", testCase.expectedValue, value)
			}
		})
	}
}

type TestStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}

func TestStructToMap(t *testing.T) {
	t.Parallel()

	obj := TestStruct{
		Field1: 10,
		Field2: "test",
		Field3: true,
	}

	result := utils.StructToMap(obj)

	expectedResult := map[string]interface{}{
		"field1": 10,
		"field2": "test",
		"field3": true,
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result: got %v, want %v", result, expectedResult)
	}
}
