package delivery

import (
	"io"
	"log"
	"net/http"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"

	"go.uber.org/zap"
)

const MaxSizePhoto = 5 * 1024 * 1024

var ErrToBigFile = myerrors.NewError("Максимальный размер фото %d Мбайт", MaxSizePhoto%1024%1024)

type IFileService interface {
	SaveImage(r io.Reader) (string, error)
}

type FileHandler struct {
	http.Handler
	fileService IFileService
	logger      *zap.SugaredLogger
}

func NewFileHandler(fileServiceDir string, fileService IFileService, logger *zap.SugaredLogger) *FileHandler {
	return &FileHandler{
		Handler:     http.StripPrefix("/api/v1/img/", http.FileServer(http.Dir(fileServiceDir))),
		fileService: fileService, logger: logger,
	}
}

// UploadFileHandler godoc
//
//	@Summary    upload photo
//	@Description  upload photo to file service and return its url
//
// @Description	@Accept      multipart-form-data TODO fix
//
//	@Tags fileService
//	@Produce    json
//	@Param      photo  body string  true  "photo row TODO fix type"
//	@Success    200  {object} ResponseURL
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /img/upload [post]
func (f *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(MaxSizePhoto)

	file, handler, err := r.FormFile("my_file")
	if err != nil {
		log.Printf("in UploadFileHandler: fileSize=%d more than max=%d", handler.Size, MaxSizePhoto)
		delivery.SendErrResponse(w, f.logger,
			delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	defer file.Close()

	if handler.Size > MaxSizePhoto {
		f.logger.Errorf("in UploadFileHandler: fileSize=%d more than max=%d", handler.Size, MaxSizePhoto)
		delivery.SendErrResponse(w, f.logger, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrToBigFile.Error()))

		return
	}

	URLToFile, err := f.fileService.SaveImage(r.Body)
	if err != nil {
		f.logger.Errorf("in UploadFileHandler: %+v\n", err)
		delivery.HandleErr(w, f.logger, err)

		return
	}

	delivery.SendOkResponse(w, f.logger, NewResponseURL(delivery.StatusResponseSuccessful, URLToFile))
	log.Printf("in UploadFileHandler: uploaded file with name=%+v", URLToFile)
}
