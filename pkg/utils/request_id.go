package utils

import (
	"context"
	"math/rand"
	"strconv"
)

var (
	requestIDKey = "req_id"

	minRequestID = 100000
	maxRequestID = 999999
)

func AddRequestIDToCtx(ctx context.Context) context.Context {
	requestID := strconv.Itoa(minRequestID + rand.Intn(maxRequestID-minRequestID+1))
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	return ctx
}

func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}
