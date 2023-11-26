package mux

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/services/file_service/internal/server/delivery"

	"go.uber.org/zap"
)

type ConfigMux struct {
	allowOrigin    string
	schema         string
	portServer     string
	fileServiceDir string
}

func NewConfigMux(allowOrigin string, schema string, portServer string, fileServiceDir string) *ConfigMux {
	return &ConfigMux{
		allowOrigin:    allowOrigin,
		schema:         schema,
		portServer:     portServer,
		fileServiceDir: fileServiceDir,
	}
}

func NewMux(ctx context.Context, configMux *ConfigMux, fileService delivery.IFileService, logger *zap.SugaredLogger,
) (http.Handler, error) {
	router := http.NewServeMux()

	fileHandler := delivery.NewFileHandler(fileService, logger, configMux.fileServiceDir)

	router.Handle("/api/v1/img/", fileHandler.DocFileServerHandler(ctx))
	router.Handle("/api/v1/img/upload", middleware.Context(ctx,
		middleware.SetupCORS(fileHandler.UploadFileHandler, configMux.allowOrigin, configMux.schema)))

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Panic(middleware.AccessLogMiddleware(router, logger), logger))

	return mux, nil
}
