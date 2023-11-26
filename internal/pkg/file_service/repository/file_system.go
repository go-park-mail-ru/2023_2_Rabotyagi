package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"log"
	"os"
)

type FileSystemStorage struct {
	baseDir string
}

func NewFileSystemStorage(baseDir string) *FileSystemStorage {
	return &FileSystemStorage{baseDir: baseDir}
}

func (f *FileSystemStorage) SaveFile(content []byte, fileName string) error {
	file, err := os.Create(f.baseDir + "/" + fileName)
	if err != nil {
		log.Printf("in Save: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Printf("in Save: %+v\n", err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
