package mylogger_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
)

func TestAddRequestIDToCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx = mylogger.AddRequestIDToCtx(ctx)

	requestID := mylogger.GetRequestIDFromCtx(ctx)

	if requestID == "" {
		t.Error("Failed to add request ID to context")
	}
}

func TestGetRequestIDFromCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	requestID := mylogger.GetRequestIDFromCtx(ctx)

	if requestID != "" {
		t.Error("Unexpected request ID in empty context")
	}

	ctx = mylogger.AddRequestIDToCtx(ctx)

	requestID = mylogger.GetRequestIDFromCtx(ctx)
	if requestID == "" {
		t.Error("Failed to add request ID to context")
	}
}
