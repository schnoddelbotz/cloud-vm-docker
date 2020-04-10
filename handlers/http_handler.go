package handlers

import "net/http"

func CloudTaskZipZap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"foo":"bar"}`))
}
