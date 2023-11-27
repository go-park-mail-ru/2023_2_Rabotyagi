package usecases

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	fileservicerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
	"go.uber.org/zap"
	"image"
	"io"
	// Add png formats for image
	_ "image/png"
	// Add jpeg formats for image
	_ "image/jpeg"
)

var (
	ErrCantRead    = myerrors.NewErrorBadContentRequest("Не получилось считать содержимое файла из тела запроса")
	ErrWrongFormat = myerrors.NewErrorBadFormatRequest("Формат файла должен быть png, jpeg")
)

var _ IFileStorageHTTP = (*fileservicerepo.FileSystemStorage)(nil)

type IFileStorageHTTP interface {
	SaveFile(content []byte, fileName string) error
}

type FileServiceHTTP struct {
	urlPrefixPath string
	fileStorage   IFileStorageHTTP
	logger        *zap.SugaredLogger
}

func NewFileServiceHTTP(fileStorage IFileStorageHTTP, urlPrefixPath string) (*FileServiceHTTP, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, myerrors.NewErrorInternal(err.Error())
	}

	return &FileServiceHTTP{fileStorage: fileStorage, urlPrefixPath: urlPrefixPath, logger: logger}, nil
}

func (f *FileServiceHTTP) SaveImage(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		f.logger.Infoln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrCantRead)
	}

	_, format, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		f.logger.Infoln(err)

		return "", fmt.Errorf("вы используете формат %s %w", format, ErrWrongFormat)
	}

	fileName, err := HashContent(content)
	if err != nil {
		f.logger.Infoln(err)

		return "", myerrors.NewErrorInternal(err.Error())
	}

	err = f.fileStorage.SaveFile(content, fileName)
	if err != nil {
		f.logger.Infoln(err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return f.urlPrefixPath + fileName, nil
}
