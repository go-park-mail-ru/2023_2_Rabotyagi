package repository

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
)

type FileSystemStorage struct {
	baseDir    string
	mapFiles   map[string]struct{}
	muMapFiles *sync.RWMutex
	logger     *my_logger.MyLogger
}

func NewFileSystemStorage(baseDir string) (*FileSystemStorage, error) {
	logger, err := my_logger.Get()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	prevFSStorage := &FileSystemStorage{
		baseDir:    baseDir,
		mapFiles:   make(map[string]struct{}),
		muMapFiles: &sync.RWMutex{},
		logger:     logger,
	}

	err = prevFSStorage.recover()
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return prevFSStorage, nil
}

func (f *FileSystemStorage) recover() error {
	files, err := os.ReadDir(f.baseDir)
	if err != nil {
		f.logger.Infoln(err)

		return myerrors.NewErrorInternal(err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			f.muMapFiles.Lock()
			f.mapFiles[file.Name()] = struct{}{}
			f.muMapFiles.Unlock()
		}
	}

	return nil
}

// Check bool in return slice means file exist if it's true.
func (f *FileSystemStorage) Check(_ context.Context, files []string) ([]bool, error) {
	result := make([]bool, len(files))

	for i, filename := range files {
		f.muMapFiles.RLock() // May be make sense block muMapFiles outside cycle
		if _, ok := f.mapFiles[filename]; ok {
			result[i] = true
		} else {
			result[i] = false
		}
		f.muMapFiles.RUnlock()
	}

	return result, nil
}

func (f *FileSystemStorage) SaveFile(ctx context.Context, content []byte, fileName string) error {
	logger := f.logger.LogReqID(ctx)

	file, err := os.Create(f.baseDir + "/" + fileName)
	if err != nil {
		logger.Infoln(err)

		return myerrors.NewErrorInternal(err.Error())
	}

	_, err = file.Write(content)
	if err != nil {
		logger.Infoln(err)

		return myerrors.NewErrorInternal(err.Error())
	}

	f.muMapFiles.Lock()
	f.mapFiles[fileName] = struct{}{}
	f.muMapFiles.Unlock()

	return nil
}
