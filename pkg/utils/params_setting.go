package utils

import (
	"net/http"
)

func AddQueryParamsToRequest(r *http.Request, params map[string]string) {
	query := r.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	r.URL.RawQuery = query.Encode()
}
