package interceptors

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/metrics"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"google.golang.org/grpc"
)

type GrpcAccessInterceptor struct {
	metrics metrics.IMetricManagerGrpc
}

func NewGrpcAccessInterceptor(metrics metrics.IMetricManagerGrpc) *GrpcAccessInterceptor {
	return &GrpcAccessInterceptor{metrics: metrics}
}

func (g *GrpcAccessInterceptor) AccessInterceptor(ctx context.Context, req interface{},
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

	g.metrics.IncTotal(info.FullMethod)
	g.metrics.AddDuration(info.FullMethod, duration)

	if errHandler != nil {
		logger.Errorln(errHandler)
		g.metrics.IncTotalErr(info.FullMethod)

		return nil, fmt.Errorf(myerrors.ErrTemplate, errHandler)
	}

	logger.Infof("method: %v request: %v duration: %v", info.FullMethod, req, duration)

	return resp, nil
}
