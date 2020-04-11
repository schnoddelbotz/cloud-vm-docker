package handlers

import "net/http"

// CloudTaskZipZap HTTP CloudFunction handler makes VMs triggerable via plain https+token request
func CloudTaskZipZap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"foo":"bar"}`))
}
