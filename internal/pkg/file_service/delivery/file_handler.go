package delivery

import (
	"log"
	"net/http"

	myerrors "github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/file_service/usecases"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/server/delivery"
)

const MaxSizePhoto = 5 * 1024 * 1024

var ErrToBigFile = myerrors.NewError("Максимальный размер фото %d Мбайт", MaxSizePhoto%1024%1024)

type FileHandler struct {
	http.Handler
	fileService usecases.IFileService
}

func NewFileHandler(baseDir string, fileService usecases.IFileService) *FileHandler {
	return &FileHandler{Handler: http.FileServer(http.Dir(baseDir)), fileService: fileService}
}

// UploadFileHandler godoc
//
//	@Summary    upload photo
//	@Description  upload photo to file service and return its url
//	@Tags product
//	@Accept      multipart-form-data TODO fix
//	@Produce    json
//	@Param      photo  body  true  "photo row"
//	@Success    200  {object} ResponseURL
//	@Failure    405  {string} string
//	@Failure    500  {string} string
//	@Failure    222  {object} delivery.ErrorResponse "Error"
//	@Router      /img/upload [post]
func (f *FileHandler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err := r.ParseMultipartForm(MaxSizePhoto)

	file, handler, err := r.FormFile("my_file")
	if err != nil {
		log.Printf("in UploadFileHandler: fileSize=%d more than max=%d", handler.Size, MaxSizePhoto)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrInternalServer, delivery.ErrInternalServer))

		return
	}

	defer file.Close()

	if handler.Size > MaxSizePhoto {
		log.Printf("in UploadFileHandler: fileSize=%d more than max=%d", handler.Size, MaxSizePhoto)
		delivery.SendErrResponse(w, delivery.NewErrResponse(delivery.StatusErrBadRequest, ErrToBigFile.Error()))

		return
	}

	URLToFile, err := f.fileService.SaveImage(r.Body)
	if err != nil {
		delivery.HandleErr(w, "in UploadFileHandler: ", err)

		return
	}

	delivery.SendOkResponse(w, NewResponseURL(delivery.StatusResponseSuccessful, URLToFile))
	log.Printf("in UploadFileHandler: uploaded file with name=%+v", URLToFile)
}
