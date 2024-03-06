package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		answer   string
		redirect string
		status   int
		size     int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write data to logger.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	r.responseData.answer = string(b)
	return size, err
}

// Log status code and add redirect header if existed.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
	if statusCode == http.StatusTemporaryRedirect {
		r.responseData.redirect = r.Header().Get("Location")
	}
}

// Middleware for logging web server (URI, Time, Method, Delay, User Info) for logging.
func MiddlwLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status:   0,
			size:     0,
			answer:   "",
			redirect: "",
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		//delay
		duration := time.Since(start)

		zap.S().Infoln(
			"Request URL from: ", uri,
			"Method: ", method,
			"Status: ", responseData.status,
			"Duration: ", duration,
			"Server answer: ", responseData.answer,
			"Size: ", responseData.size,
		)
		if responseData.status == http.StatusTemporaryRedirect {
			zap.S().Infoln("Redirect to: ", responseData.redirect)
		}
	})

}
