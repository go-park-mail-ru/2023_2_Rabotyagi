package delivery

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/usecases"
)

const (
	MaxSizePhotoBytes = 5 * 1024 * 1024
	MaxCountPhoto     = 4

	nameImagesInForm = "images"

	rootPath = "/api/v1/img/"
)

type keyCtx string

const keyCtxHandler keyCtx = "handler"

var (
	ErrToBigFile = myerrors.NewErrorBadContentRequest("Максимальный размер фото %d Мбайт",
		MaxSizePhotoBytes%1024%1024)
	ErrToManyCountFiles  = myerrors.NewErrorBadContentRequest("Максимальное количество фото = %d", MaxCountPhoto)
	ErrForbiddenRootPath = myerrors.NewErrorBadContentRequest("Нельзя вызывать корневой путь")
)

var _ IFileServiceHTTP = (*fileusecases.FileServiceHTTP)(nil)

type IFileServiceHTTP interface {
	SaveImage(r io.Reader) (string, error)
}

type FileHandlerHTTP struct {
	fileServiceDir string
	fileService    IFileServiceHTTP
	logger         *my_logger.MyLogger
}

func NewFileHandlerHTTP(fileService IFileServiceHTTP,
	logger *my_logger.MyLogger, fileServiceDir string,
) *FileHandlerHTTP {
	return &FileHandlerHTTP{
		fileService: fileService, logger: logger,
		fileServiceDir: fileServiceDir,
	}
}

// UploadFileHandler godoc
//
//	@Summary    upload photo
//	@Description  upload photo to file service and return its url
//
//	@Tags fileService
//
//	@Accept     multipart/form-data
//	@Produce    json
//	@Success    200  {object} ResponseURLs
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Тут статус http статус 200. Внутри body статус может быть badContent(4400), badFormat(4000)"
//	@Router      /img/upload [post]
func (f *FileHandlerHTTP) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxSizePhotoBytes*MaxCountPhoto)
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	logger := f.logger.LogReqID(r.Context())

	err := r.ParseMultipartForm(MaxSizePhotoBytes)
	if err != nil {
		logger.Errorln(err)
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	slFiles, ok := r.MultipartForm.File[nameImagesInForm]
	if !ok {
		logger.Errorln(err)
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	if len(slFiles) > MaxCountPhoto {
		logger.Errorln(ErrToManyCountFiles)
		responses.HandleErr(w, logger, ErrToManyCountFiles)

		return
	}

	slURL := make([]string, len(slFiles))

	for i, file := range slFiles {
		if file.Size > MaxSizePhotoBytes {
			err := myerrors.NewErrorBadContentRequest(
				"файл: %s весит %d Мбайт. %+v\n", file.Filename, file.Size/1024/1024, ErrToBigFile.Error())

			logger.Errorln(err)
			responses.HandleErr(w, logger, err)

			return
		}

		fileBody, err := file.Open()
		if err != nil {
			logger.Errorln(err)
			responses.SendResponse(w, logger,
				responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

			return
		}

		URLToFile, err := f.fileService.SaveImage(fileBody)
		if err != nil {
			logger.Errorln(err)
			responses.HandleErr(w, logger, err)

			return
		}

		slURL[i] = URLToFile
	}

	responses.SendResponse(w, logger, NewResponseURLs(slURL))

	for _, fileName := range slURL {
		logger.Infof("uploaded file %s", fileName)
	}
}

// fileServerHandler godoc
//
//	@Summary    download photo
//	@Description  download photo of file by its name
//
//	@Tags fileService
//
//	@Accept     json
//	@Produce    png
//	@Produce    jpeg
//	@Produce    json
//	@Param      name path  string true "name of image"
//	@Failure    405  {string} string
//	@Failure    404  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} responses.ErrorResponse "Тут статус http статус 200. Внутри body статус может быть badContent(4400)"
//	@Router      /img/ [get]
func (f *FileHandlerHTTP) fileServerHandler(w http.ResponseWriter, r *http.Request) {
	logger := f.logger.LogReqID(r.Context())

	if r.URL.Path == rootPath {
		logger.Errorln(ErrForbiddenRootPath)
		responses.HandleErr(w, logger, ErrForbiddenRootPath)

		return
	}

	ctx := r.Context()
	fileServerRaw := ctx.Value(keyCtxHandler)

	fileServer, ok := fileServerRaw.(http.Handler)
	if !ok {
		logger.Errorln(fmt.Sprintf("handler = %+v а должен быть типом http.Handler", fileServerRaw))
		responses.SendResponse(w, logger,
			responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer))

		return
	}

	fileServer.ServeHTTP(w, r)
}

func (f *FileHandlerHTTP) DocFileServerHandler() http.Handler {
	fileServer := http.StripPrefix("/api/v1/img/", http.FileServer(http.Dir(f.fileServiceDir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), keyCtxHandler, fileServer))

		f.fileServerHandler(w, r)
	})
}