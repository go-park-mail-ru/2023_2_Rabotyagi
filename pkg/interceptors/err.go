package interceptors

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrConvertInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)

	myErr := &myerrors.Error{}
	if errors.As(err, &myErr) && myErr.IsErrorClient() {
		err = status.Errorf(codes.Code(myErr.Status()), myErr.Error())

		return resp, err
	}

	return resp, err
}
