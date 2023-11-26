package usecases

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	fileservicerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
	"image"
	"io"
	"log"

	// Add png formats for image
	_ "image/png"
	// Add jpeg formats for image
	_ "image/jpeg"
)

var (
	ErrCantRead    = myerrors.NewErrorInternal("Не получилось считать содержимое файла из тела запроса")
	ErrWrongFormat = myerrors.NewErrorBadFormatRequest("Формат файла должен быть png, jpeg")
)

var _ IFileStorage = (*fileservicerepo.FileSystemStorage)(nil)

type IFileStorage interface {
	SaveFile(content []byte, fileName string) error
}

type FileService struct {
	urlPrefixPath string
	fileStorage   IFileStorage
}

func NewFileService(fileStorage IFileStorage, urlPrefixPath string) *FileService {
	return &FileService{fileStorage: fileStorage, urlPrefixPath: urlPrefixPath}
}

func (f *FileService) SaveImage(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("in SaveImage: %+v\n", err)

		return "", fmt.Errorf(myerrors.ErrTemplate, ErrCantRead)
	}

	_, format, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		log.Printf("in SaveImage: %+v\n", err)

		return "", fmt.Errorf("вы используете формат %s %w", format, ErrWrongFormat)
	}

	fileName, err := HashContent(content)
	if err != nil {
		log.Printf("in SaveImage: %+v\n", err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = f.fileStorage.SaveFile(content, fileName)
	if err != nil {
		log.Printf("in SaveImage: %+v\n", err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return f.urlPrefixPath + fileName, nil
}
