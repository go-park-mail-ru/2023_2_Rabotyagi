package middleware

import (
	"net/http"
)

func setupHeadersCORS(w http.ResponseWriter, allowOrigin string, schema string) {
	w.Header().Set("Access-Control-Allow-Origin", schema+allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func SetupCORS(next http.HandlerFunc, addrOrigin string, schema string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupHeadersCORS(w, addrOrigin, schema)

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
