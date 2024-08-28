package main

import "net/http"

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
