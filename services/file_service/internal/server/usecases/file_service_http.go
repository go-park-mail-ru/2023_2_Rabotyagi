package usecases

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	// Add png formats for image
	_ "image/png"
	// Add jpeg formats for image
	_ "image/jpeg"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	fileservicerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
)

var (
	ErrCantRead    = myerrors.NewErrorBadFormatRequest("Не получилось считать содержимое файла из тела запроса")
	ErrWrongFormat = myerrors.NewErrorBadContentRequest("Формат файла должен быть png, jpeg")
)

var _ IFileStorageHTTP = (*fileservicerepo.FileSystemStorage)(nil)

type IFileStorageHTTP interface {
	SaveFile(ctx context.Context, content []byte, fileName string) error
}

type FileServiceHTTP struct {
	urlPrefixPath string
	fileStorage   IFileStorageHTTP
	logger        *my_logger.MyLogger
}

func NewFileServiceHTTP(fileStorage IFileStorageHTTP, urlPrefixPath string) (*FileServiceHTTP, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return &FileServiceHTTP{fileStorage: fileStorage, urlPrefixPath: urlPrefixPath, logger: logger}, nil
}

func (f *FileServiceHTTP) SaveImage(ctx context.Context, reader io.Reader) (string, error) {
	logger := f.logger.LogReqID(ctx)

	content, err := io.ReadAll(reader)
	if err != nil {
		logger.Infoln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrCantRead)
	}

	_, format, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		logger.Infoln(err)

		return "", fmt.Errorf("вы используете формат %s %w", format, ErrWrongFormat)
	}

	fileName, err := HashContent(ctx, content)
	if err != nil {
		logger.Infoln(err)

		return "", myerrors.NewErrorInternal(err.Error())
	}

	err = f.fileStorage.SaveFile(ctx, content, fileName)
	if err != nil {
		logger.Infoln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return f.urlPrefixPath + fileName, nil
}
