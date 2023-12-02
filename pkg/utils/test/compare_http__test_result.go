package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
)

func CompareSameType[T comparable](received T, expected T) error {
	if received != expected {
		return fmt.Errorf("response: got %d, expected %d",
			received, expected)
	}

	return nil
}

func CompareHTTPTestResult(recorder *httptest.ResponseRecorder, expected any) error {
	resp := recorder.Result()
	defer resp.Body.Close()

	receivedRespRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to ReadAll resp.Body: %w", err)
	}

	if expectedStr, ok := expected.(string); ok {
		receivedStr := string(receivedRespRaw)
		if expectedStr != receivedStr {
			return fmt.Errorf("response: got %s, expected %s",
				receivedStr, expectedStr)
		}

		return nil
	}

	expectedResponseRaw, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal expexted: %w", err)
	}

	if !bytes.Equal(receivedRespRaw, expectedResponseRaw) {
		return fmt.Errorf("response: got %s, expected %s",
			string(receivedRespRaw), string(expectedResponseRaw))
	}

	return nil
}
