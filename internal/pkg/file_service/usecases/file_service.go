package usecases

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"

	// Add png formats for image
	_ "image/png"
	// Add jpeg formats for image
	_ "image/jpeg"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	fileservicerepo "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/repository"
)

var _ IFileStorage = (*fileservicerepo.FileSystemStorage)(nil)

var (
	ErrCantRead    = myerrors.NewError("Не получилось считать содержимое файла")
	ErrWrongFormat = myerrors.NewError("Формат файла должен быть png, jpeg")
)

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

		return "", fmt.Errorf("%w", ErrWrongFormat)
	}

	fileName := HashContent(content) + format

	err = f.fileStorage.SaveFile(content, fileName)
	if err != nil {
		log.Printf("in SaveImage: %+v\n", err)

		return "", fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return f.urlPrefixPath + fileName, nil
}
