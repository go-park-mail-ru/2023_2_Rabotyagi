package my_logger

import (
	"context"
	"math/rand"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type keyCtx string

const (
	requestIDKey keyCtx = "req_id"

	minRequestID = 100000
	maxRequestID = 999999
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	return ctx
}

func AddRequestIDToCtx(ctx context.Context) context.Context {
	requestID := strconv.Itoa(minRequestID + rand.Intn(maxRequestID-minRequestID+1)) //nolint:gosec

	return SetRequestIDToCtx(ctx, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

func NewMDFromRequestIDCtx(ctx context.Context) metadata.MD {
	return metadata.Pairs(string(requestIDKey), GetRequestIDFromCtx(ctx))
}

func GetRequestIDFromMDCtx(ctx context.Context) string {
	slStr := metadata.ValueFromIncomingContext(ctx, string(requestIDKey))
	if len(slStr) < 1 {
		return ""
	}

	return slStr[0]
}
