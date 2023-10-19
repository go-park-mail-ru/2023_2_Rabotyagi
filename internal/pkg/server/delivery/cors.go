package delivery

import "net/http"

// Schema TODO from config.
const Schema = "http://"

func SetupCORS(w http.ResponseWriter, allowOrigin string) {
	w.Header().Set("Access-Control-Allow-Origin", Schema+allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
