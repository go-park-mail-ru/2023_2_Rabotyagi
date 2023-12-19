package middleware

import (
	"net/http"
	"strings"
)

const APINameV1 = "/api/v1"

func AddAPIName(handler http.Handler, apiName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPath := strings.TrimPrefix(r.URL.Path, apiName)
		if newPath == r.URL.Path {
			http.Error(w, `Не найден`, http.StatusNotFound)

			return
		}

		r.URL.Path = newPath
		handler.ServeHTTP(w, r)
	})
}
