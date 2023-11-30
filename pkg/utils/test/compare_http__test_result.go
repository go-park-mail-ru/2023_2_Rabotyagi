package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
)

func CompareHTTPTestResult(recorder *httptest.ResponseRecorder, expected any) error {
	resp := recorder.Result()
	defer resp.Body.Close()

	receivedRespRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to ReadAll resp.Body: %w", err)
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
