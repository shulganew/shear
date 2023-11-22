package handlers

import (
	"net/http"
	"time"

	"github.com/shulganew/shear.git/internal/config"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// middleware for logging web server
// URI
// Time
// Method
// Delay
// User Info for logging
func MidlewLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		//delay
		duration := time.Since(start)
		logz := config.InitLog()
		logz.Infoln(
			"URI: ", uri,
			"Method: ", method,
			"Status: ", responseData.status,
			"Duration: ", duration,
			"Size: ", responseData.size,
		)
	})

}
