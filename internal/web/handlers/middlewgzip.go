package handlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shulganew/shear.git/internal/api"
	"github.com/shulganew/shear.git/internal/config"
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
		logz := config.InitLog()
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		//check if client send gzip json

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			logz.Infoln("Client send a file in gzip format!")
			var reader io.Reader

			if r.Header.Get(`Content-Encoding`) == `gzip` {
				gz, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				reader = gz
				defer gz.Close()
			} else {
				reader = r.Body
			}

			body, err := io.ReadAll(reader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//update body with unzipped file
			read := bytes.NewReader(body)
			readCloser := io.NopCloser(read)
			r.Header.Del("Content-Encoding")
			r.Header.Set("Content-Encoding", api.CONTENT_TYPE_JSON)
			r.Body = readCloser

		}

		//if browser doesnt support gzip exit
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			logz.Infoln("Clinent does not support gzip format!")
			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			logz.Errorln("error during gzip compression")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		logz.Infoln("Send zipped answer!!!!!!!!!!!!!!!!!!!")
		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

		duration := time.Since(start)

		logz.Infoln(
			"URI: ", uri,
			"Method: ", method,
			"Duration zip json: ", duration,
		)
	})

}
