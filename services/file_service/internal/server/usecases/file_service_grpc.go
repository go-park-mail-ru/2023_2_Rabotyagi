package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
)

var _ IFileStorageGrpc = (*repository.FileSystemStorage)(nil)

type IFileStorageGrpc interface {
	Check(ctx context.Context, files []string) ([]bool, error)
}

type FileServiceGrpc struct {
	urlPrefixPath string
	fileStorage   IFileStorageGrpc
}

func NewFileServiceGrpc(urlPrefixPath string, fileStorage IFileStorageGrpc) *FileServiceGrpc {
	return &FileServiceGrpc{urlPrefixPath: urlPrefixPath, fileStorage: fileStorage}
}

func (f *FileServiceGrpc) Check(ctx context.Context, urls []string) ([]bool, error) {
	for i, url := range urls {
		urls[i] = strings.TrimPrefix(url, f.urlPrefixPath)
	}

	result, err := f.fileStorage.Check(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return result, nil
}
