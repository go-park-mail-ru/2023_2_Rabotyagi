package my_logger_test

import (
	"context"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"testing"
)

func TestAddRequestIDToCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx = my_logger.AddRequestIDToCtx(ctx)

	requestID := my_logger.GetRequestIDFromCtx(ctx)

	if requestID == "" {
		t.Error("Failed to add request ID to context")
	}
}

func TestGetRequestIDFromCtx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	requestID := my_logger.GetRequestIDFromCtx(ctx)

	if requestID != "" {
		t.Error("Unexpected request ID in empty context")
	}

	ctx = my_logger.AddRequestIDToCtx(ctx)

	requestID = my_logger.GetRequestIDFromCtx(ctx)
	if requestID == "" {
		t.Error("Failed to add request ID to context")
	}
}
