package interceptors

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"google.golang.org/grpc"
)

func AccessInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	logger.Infof("Received request: %v", req)

	return resp, nil
}
