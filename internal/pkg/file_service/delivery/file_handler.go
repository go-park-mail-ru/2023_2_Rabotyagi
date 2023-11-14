package delivery

import (
	fileusecases "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/usecases"
	"io"
	"net/http"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"

	"go.uber.org/zap"
)

const (
	MaxSizePhotoBytes = 5 * 1024 * 1024
	MaxCountPhoto     = 4

	nameImagesInForm = "images"

	rootPath = "/api/v1/img/"
)

var (
	ErrToBigFile         = myerrors.NewError("Максимальный размер фото %d Мбайт", MaxSizePhotoBytes%1024%1024)
	ErrToManyCountFiles  = myerrors.NewError("Максимальное количество фото = %d", MaxCountPhoto)
	ErrForbiddenRootPath = myerrors.NewError("Нельзя вызывать корневой путь")
)

var _ IFileService = (*fileusecases.FileService)(nil)

type IFileService interface {
	SaveImage(r io.Reader) (string, error)
}

type FileHandler struct {
	fileServiceDir string
	addrOrigin     string
	schema         string
	fileService    IFileService
	logger         *zap.SugaredLogger
}

func NewFileHandler(fileService IFileService, logger *zap.SugaredLogger,
	fileServiceDir string, addrOrigin string, schema string,
) *FileHandler {
	return &FileHandler{
		fileService: fileService, logger: logger,
		fileServiceDir: fileServiceDir,
		addrOrigin:     addrOrigin,
		schema:         schema,
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
	delivery.SetupCORS(w, f.addrOrigin, f.schema)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `Method not allowed`, http.StatusMethodNotAllowed)

		return
	}

	err := r.ParseMultipartForm(MaxSizePhotoBytes)
	if err != nil {
		f.logger.Errorln(err)
		delivery.SendErrResponse(w, f.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	slFiles, ok := r.MultipartForm.File[nameImagesInForm]
	if !ok {
		f.logger.Errorln(err)
		delivery.SendErrResponse(w, f.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	if len(slFiles) > MaxCountPhoto {
		f.logger.Errorln(ErrToManyCountFiles)
		delivery.SendErrResponse(w, f.logger,
			delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrToManyCountFiles.Error()))

		return
	}

	slURL := make([]string, len(slFiles))

	for i, file := range slFiles {
		if file.Size > MaxSizePhotoBytes {
			f.logger.Errorf("filename = %s error: %+v\n", file.Filename, ErrToBigFile)
			delivery.SendErrResponse(w, f.logger,
				delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrToBigFile.Error()))

			return
		}

		fileBody, err := file.Open()
		if err != nil {
			f.logger.Errorln(err)
			delivery.SendErrResponse(w, f.logger,
				delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

			return
		}

		URLToFile, err := f.fileService.SaveImage(fileBody)
		if err != nil {
			f.logger.Errorln(err)
			delivery.HandleErr(w, f.logger, err)

			return
		}

		slURL[i] = URLToFile
	}

	delivery.SendOkResponse(w, f.logger, NewResponseURLs(delivery.StatusResponseSuccessful, slURL))

	for _, fileName := range slURL {
		f.logger.Infof("uploaded file %s", fileName)
	}
}

// DocHandlerFileServer godoc
// TODO not worked documentation
func (f *FileHandler) DocHandlerFileServer() http.Handler {
	fileServer := http.StripPrefix("/api/v1/img/", http.FileServer(http.Dir(f.fileServiceDir)))

	sanitizedHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == rootPath {
			f.logger.Errorln(ErrForbiddenRootPath)
			delivery.HandleErr(w, f.logger, ErrForbiddenRootPath)

			return
		}

		fileServer.ServeHTTP(w, r)
	}

	return http.HandlerFunc(sanitizedHandlerFunc)
}
