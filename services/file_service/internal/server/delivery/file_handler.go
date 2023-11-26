package delivery

import (
	"context"
	"fmt"
	delivery2 "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/server/delivery/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/usecases"
	"io"
	"net/http"

	"go.uber.org/zap"
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

var _ IFileService = (*fileusecases.FileService)(nil)

type IFileService interface {
	SaveImage(r io.Reader) (string, error)
}

type FileHandler struct {
	fileServiceDir string
	fileService    IFileService
	logger         *zap.SugaredLogger
}

func NewFileHandler(fileService IFileService,
	logger *zap.SugaredLogger, fileServiceDir string,
) *FileHandler {
	return &FileHandler{
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
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /img/upload [post]
func (f *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxSizePhotoBytes*MaxCountPhoto)
	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	err := r.ParseMultipartForm(MaxSizePhotoBytes)
	if err != nil {
		f.logger.Errorln(err)
		delivery2.SendResponse(w, f.logger,
			delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

		return
	}

	slFiles, ok := r.MultipartForm.File[nameImagesInForm]
	if !ok {
		f.logger.Errorln(err)
		delivery2.SendResponse(w, f.logger,
			delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

		return
	}

	if len(slFiles) > MaxCountPhoto {
		f.logger.Errorln(ErrToManyCountFiles)
		delivery2.SendResponse(w, f.logger,
			delivery2.NewErrResponse(statuses.StatusBadContentRequest, ErrToManyCountFiles.Error()))

		return
	}

	slURL := make([]string, len(slFiles))

	for i, file := range slFiles {
		if file.Size > MaxSizePhotoBytes {
			f.logger.Errorf("filename = %s error: %+v\n", file.Filename, ErrToBigFile)
			delivery2.SendResponse(w, f.logger,
				delivery2.NewErrResponse(statuses.StatusBadContentRequest, ErrToBigFile.Error()))

			return
		}

		fileBody, err := file.Open()
		if err != nil {
			f.logger.Errorln(err)
			delivery2.SendResponse(w, f.logger,
				delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

			return
		}

		URLToFile, err := f.fileService.SaveImage(fileBody)
		if err != nil {
			f.logger.Errorln(err)
			delivery2.HandleErr(w, f.logger, err)

			return
		}

		slURL[i] = URLToFile
	}

	delivery2.SendResponse(w, f.logger, NewResponseURLs(slURL))

	for _, fileName := range slURL {
		f.logger.Infof("uploaded file %s", fileName)
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
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /img/ [get]
func (f *FileHandler) fileServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == rootPath {
		f.logger.Errorln(ErrForbiddenRootPath)
		delivery2.HandleErr(w, f.logger, ErrForbiddenRootPath)

		return
	}

	ctx := r.Context()
	fileServerRaw := ctx.Value(keyCtxHandler)

	fileServer, ok := fileServerRaw.(http.Handler)
	if !ok {
		f.logger.Errorln(fmt.Sprintf("handler = %+v а должен быть типом http.Handler", fileServerRaw))
		delivery2.SendResponse(w, f.logger,
			delivery2.NewErrResponse(statuses.StatusInternalServer, delivery2.ErrInternalServer))

		return
	}

	fileServer.ServeHTTP(w, r)
}

func (f *FileHandler) DocFileServerHandler(ctx context.Context) http.Handler {
	fileServer := http.StripPrefix("/api/v1/img/", http.FileServer(http.Dir(f.fileServiceDir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(ctx, keyCtxHandler, fileServer))

		f.fileServerHandler(w, r)
	})
}
