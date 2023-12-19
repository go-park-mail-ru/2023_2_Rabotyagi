package delivery_test

import (
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/myerrors"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/responses/statuses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/mocks"

	"go.uber.org/mock/gomock"
)

func NewFileHandlerHTTP(ctrl *gomock.Controller,
	behaviorSessionManagerClient func(m *mocks.MockIFileServiceHTTP),
) *delivery.FileHandlerHTTP {
	mockFileServiceHTTP := mocks.NewMockIFileServiceHTTP(ctrl)

	behaviorSessionManagerClient(mockFileServiceHTTP)

	fileHandler := delivery.NewFileHandlerHTTP(mockFileServiceHTTP, my_logger.NewNop(), ".")

	return fileHandler
}

//nolint:funlen
func TestUploadFile(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name                    string
		behaviorFileServiceHTTP func(m *mocks.MockIFileServiceHTTP)
		request                 *http.Request
		expectedResponse        any
	}

	testCases := [...]TestCase{
		{
			name:                    "method not allowed",
			request:                 httptest.NewRequest(http.MethodGet, "/img/upload", nil),
			behaviorFileServiceHTTP: func(m *mocks.MockIFileServiceHTTP) {},
			expectedResponse:        "Method not allowed\n",
		},
		{
			name:                    "wrong content type",
			request:                 httptest.NewRequest(http.MethodPost, "/img/upload", nil),
			behaviorFileServiceHTTP: func(m *mocks.MockIFileServiceHTTP) {},
			expectedResponse: responses.NewErrResponse(statuses.StatusBadFormatRequest,
				"request Content-Type isn't multipart/form-data"),
		},
		{
			name: "upload test image",
			request: func() *http.Request {
				pipeReader, pipeWriter := io.Pipe()
				formWriter := multipart.NewWriter(pipeWriter)

				go func() {
					defer formWriter.Close()
					part, err := formWriter.CreateFormFile(delivery.NameImagesInForm, "test.png")
					if err != nil {
						t.Error(err)
					}

					img := image.NewNRGBA(image.Rect(0, 0, 10, 10))

					err = png.Encode(part, img)
					if err != nil {
						t.Error(err)
					}
				}()

				req := httptest.NewRequest(http.MethodPost, "/img/upload", pipeReader)
				req.Header.Set("Content-Type", formWriter.FormDataContentType())

				return req
			}(),
			behaviorFileServiceHTTP: func(m *mocks.MockIFileServiceHTTP) {
				m.EXPECT().SaveImage(gomock.Any(), gomock.Not(nil)).Return("test_url", nil)
			},
			expectedResponse: delivery.NewResponseURLs([]string{"test_url"}),
		},
		{
			name: "internal error",
			request: func() *http.Request {
				pipeReader, pipeWriter := io.Pipe()
				formWriter := multipart.NewWriter(pipeWriter)

				go func() {
					defer formWriter.Close()
					part, err := formWriter.CreateFormFile(delivery.NameImagesInForm, "test.png")
					if err != nil {
						t.Error(err)
					}

					img := image.NewNRGBA(image.Rect(0, 0, 10, 10))

					err = png.Encode(part, img)
					if err != nil {
						t.Error(err)
					}
				}()

				req := httptest.NewRequest(http.MethodPost, "/img/upload", pipeReader)
				req.Header.Set("Content-Type", formWriter.FormDataContentType())

				return req
			}(),
			behaviorFileServiceHTTP: func(m *mocks.MockIFileServiceHTTP) {
				m.EXPECT().SaveImage(gomock.Any(), gomock.Not(nil)).Return("", myerrors.NewErrorInternal("Test err"))
			},
			expectedResponse: responses.NewErrResponse(statuses.StatusInternalServer, responses.ErrInternalServer),
		},
		{
			name: "form with wrong field-name",
			request: func() *http.Request {
				pipeReader, pipeWriter := io.Pipe()
				formWriter := multipart.NewWriter(pipeWriter)

				go func() {
					defer formWriter.Close()
					part, err := formWriter.CreateFormFile("wrong field-name", "test.png")
					if err != nil {
						t.Error(err)
					}

					img := image.NewNRGBA(image.Rect(0, 0, 10, 10))

					err = png.Encode(part, img)
					if err != nil {
						t.Error(err)
					}
				}()

				req := httptest.NewRequest(http.MethodPost, "/img/upload", pipeReader)
				req.Header.Set("Content-Type", formWriter.FormDataContentType())

				return req
			}(),
			behaviorFileServiceHTTP: func(m *mocks.MockIFileServiceHTTP) {},
			expectedResponse: responses.NewErrResponse(delivery.ErrWrongNameMultipart.Status(),
				delivery.ErrWrongNameMultipart.Error()),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileHandler := NewFileHandlerHTTP(ctrl, testCase.behaviorFileServiceHTTP)

			w := httptest.NewRecorder()

			profileHandler.UploadFileHandler(w, testCase.request)

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

//nolint:nolintlint,funlen
func TestDocFileHandler(t *testing.T) {
	t.Parallel()

	_ = my_logger.NewNop()

	type TestCase struct {
		name             string
		request          *http.Request
		expectedResponse any
	}

	testCases := [...]TestCase{
		{
			name:             "test not found",
			request:          httptest.NewRequest(http.MethodGet, "/api/v1/img/file_for_test.txt", nil),
			expectedResponse: "404 page not found\n",
		},
		{
			name:             "method not allowed",
			request:          httptest.NewRequest(http.MethodDelete, "/api/v1/img/file_for_test.txt", nil),
			expectedResponse: "Method not allowed\n",
		},
		{
			name:    "can`t request root",
			request: httptest.NewRequest(http.MethodGet, "/api/v1/img/", nil),
			expectedResponse: responses.NewErrResponse(delivery.ErrForbiddenRootPath.Status(),
				delivery.ErrForbiddenRootPath.Error()),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			profileHandler := NewFileHandlerHTTP(ctrl, func(m *mocks.MockIFileServiceHTTP) {})
			docFileServer := profileHandler.DocFileServerHandler()

			w := httptest.NewRecorder()

			docFileServer.ServeHTTP(w, testCase.request)

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}
