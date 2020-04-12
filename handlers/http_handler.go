package handlers

import "net/http"

// CloudTaskZipZap HTTP CloudFunction handler makes VMs triggerable via plain https+token request
func CloudTaskZipZap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"foo":"bar"}`))
	// TODO:
	// This should become an HTTP entrypoint for exposing ctzz functionality via simple JSON api.
	// While using ctzz is obviously the most simple/direct approach to submit tasks or
	// manage them -as it speaks to google services like PubSub and FireStore directly-
	// it may be helpful to have a RESTish entrypoint for lightweight submission
	// scenarios relying entirely on e.g. just curl.
	// Obviously, it should be (auto-generated-if-not-provided) token-protected, or
	// if the user chooses so, only callable with valid IAM credentials (non-public http endpoint).
}
