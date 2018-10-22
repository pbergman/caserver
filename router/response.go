package router

import (
	"net/http"
	"strconv"
)

// wrapper around the ResponseWriter so we can
// catch the status code for debugging/logging
type wrappedResponse struct {
	http.ResponseWriter
	statusCode int
}

func newWrappedResponse(inner http.ResponseWriter) http.ResponseWriter {
	return &wrappedResponse{inner, http.StatusOK}
}

func (r *wrappedResponse) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r wrappedResponse) statusLine() string {
	return strconv.Itoa(r.statusCode) + " " + http.StatusText(r.statusCode)
}
