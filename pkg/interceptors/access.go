package interceptors

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"google.golang.org/grpc"
)

func AccessInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	reqID := my_logger.GetRequestIDFromMDCtx(ctx)
	ctx = my_logger.SetRequestIDToCtx(ctx, reqID)

	start := time.Now()
	resp, errHandler := handler(ctx, req)
	duration := time.Since(start)

	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	if errHandler != nil {
		logger.Errorln(errHandler)

		return nil, fmt.Errorf(myerrors.ErrTemplate, errHandler)
	}

	logger.Infof("method: %v request: %v duration: %v", info.FullMethod, req, duration)

	return resp, nil
}
