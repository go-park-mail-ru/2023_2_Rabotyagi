package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func AddQueryParamsToRequest(r *http.Request, params map[string]string) {
	query := r.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}

	r.URL.RawQuery = query.Encode()
}

func AddJSONBodyToRequest(r *http.Request, jsonStr string) {
	r.Header.Set("Content-Type", "application/json")
	r.Body = ioutil.NopCloser(strings.NewReader(jsonStr))
	r.ContentLength = int64(len(jsonStr))
}
