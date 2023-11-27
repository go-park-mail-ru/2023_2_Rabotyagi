package usecases

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/repository"
)

var _ IFileStorageGrpc = (*repository.FileSystemStorage)(nil)

type IFileStorageGrpc interface {
	Check(ctx context.Context, files []string) ([]bool, error)
}

type FileServiceGrpc struct {
	fileStorage IFileStorageGrpc
}

func NewFileServiceGrpc(fileStorage IFileStorageGrpc) *FileServiceGrpc {
	return &FileServiceGrpc{fileStorage: fileStorage}
}

func (f *FileServiceGrpc) Check(ctx context.Context, files []string) ([]bool, error) {
	result, err := f.fileStorage.Check(ctx, files)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return result, nil
}
