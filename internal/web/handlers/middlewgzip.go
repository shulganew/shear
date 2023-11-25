package handlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shulganew/shear.git/internal/appconsts"
	"go.uber.org/zap"
)

type gzipRequest struct {
	Req *http.Request
}

func (r gzipRequest) newBody(body io.ReadCloser) {
	r.Req.Body = body
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func MidlewZip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		if !strings.Contains(r.Header.Get(appconsts.AcceptEncoding), "gzip") {
			zap.S().Infoln("Clinent does not support gzip format!")
			h.ServeHTTP(w, r)
			return
		}

		//check if client send compressed content in the body (gzip only)
		if strings.Contains(r.Header.Get(appconsts.ContentEncoding), "gzip") {
			zap.S().Infoln("Client send a file in gzip format!")

			var reader io.Reader
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				zap.S().Errorln("Error unzip reques body")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()

			body, err := io.ReadAll(reader)
			if err != nil {
				zap.S().Errorln("Error read data from unzipped reques body")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//update body with unzipped file
			read := bytes.NewReader(body)
			readCloser := io.NopCloser(read)

			//send to ServeHTTP without encoding
			r.Header.Del(appconsts.ContentEncoding)
			r.Body = readCloser
		}

		//r.Header.Set("Content-Type", appconsts.ContentTypeJSON)
		//r.Header.Set(appconsts.ContentEncoding, appconsts.ContentTypeJSON)

		//Send compressed with gzip unsver
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			zap.S().Errorln("error during gzip compression")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set(appconsts.ContentEncoding, "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

		duration := time.Since(start)

		zap.S().Infoln(
			"URI: ", uri,
			"Method: ", method,
			"Duration zip json: ", duration,
		)
	})

}
