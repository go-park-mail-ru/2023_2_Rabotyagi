package delivery

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"

	fileservice "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/file_service"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/usecases"
)

var ErrImgUrlsNil = myerrors.NewErrorInternal("imgURLs == nil")

var _ IFileServiceGrpc = (*fileusecases.FileServiceGrpc)(nil)

type IFileServiceGrpc interface {
	Check(ctx context.Context, files []string) ([]bool, error)
}

type FileHandlerGrpc struct {
	fileservice.UnimplementedFileServiceServer

	fileService IFileServiceGrpc
}

func NewFileHandlerGrpc(fileService IFileServiceGrpc) *FileHandlerGrpc {
	return &FileHandlerGrpc{fileService: fileService} //nolint:exhaustruct
}

func (f *FileHandlerGrpc) Check(ctx context.Context, imgURLs *fileservice.ImgURLs) (*fileservice.CheckedURLs, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	logger = logger.LogReqID(ctx)

	if imgURLs == nil {
		logger.Errorln(ErrImgUrlsNil)

		return nil, ErrImgUrlsNil
	}

	result, err := f.fileService.Check(ctx, imgURLs.GetUrl())
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &fileservice.CheckedURLs{Correct: result}, nil
}
