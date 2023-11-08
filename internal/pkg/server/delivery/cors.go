package delivery

import "net/http"

func SetupCORS(w http.ResponseWriter, allowOrigin string, schema string) {
	w.Header().Set("Access-Control-Allow-Origin", schema+allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
