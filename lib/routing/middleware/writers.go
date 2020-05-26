package middleware

import "net/http"

type statusResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

var _ http.ResponseWriter = (*statusResponseWriter)(nil)

func (sw *statusResponseWriter) Header() http.Header { return sw.w.Header() }

func (sw *statusResponseWriter) Write(b []byte) (int, error) { return sw.w.Write(b) }

func (sw *statusResponseWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.w.WriteHeader(statusCode)
}
