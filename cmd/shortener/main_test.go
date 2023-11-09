package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	webhandl "github.com/shulganew/shear.git/internal/handlers"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// get shortUrl from BDUrl
func getShortUrl(m map[string]string, longUrl string) (shortUrl string, ok bool) {
	for k, v := range m {
		if v == longUrl {
			shortUrl = k
			ok = true
			return
		}
	}
	return
}

func Test_main(t *testing.T) {
	tests := []struct {
		name        string
		request     string
		body        string
		method      string
		contentType string
		statusCode  int
		//want
	}{
		{
			name:        "base test POTS",
			request:     "http://localhost:8080",
			body:        "http://yandex.ru/",
			method:      http.MethodPost,
			contentType: "text/plain",
			statusCode:  201,
		},

		{
			name:        "base test GET",
			request:     "http://localhost:8080",
			body:        "http://yandex.ru/",
			method:      http.MethodGet,
			contentType: "text/plain",
			statusCode:  307,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.method == http.MethodPost {
				t.Log("=============POTS===============")
				t.Log("tt.request=", tt.request)
				t.Log("strings.NewReader(tt.body)=", tt.body)
				request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
				//create status recorder
				resRecord := httptest.NewRecorder()

				webhandl.SetUrl(resRecord, request)

				//get result
				res := resRecord.Result()

				//check answer code
				t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
				assert.Equal(t, tt.statusCode, res.StatusCode)

				//check Content type
				t.Log("Content-Type exp: ", tt.contentType, " act: ", res.Header.Get("Content-Type"))
				assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))

				//check body content
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				body := string(resBody)
				t.Log("Body: ", body)

				fullUrl, err := url.Parse(body)
				if err != nil {
					panic(err)
				}
				t.Log("full url: ", fullUrl.Path)
				shortUrl := strings.TrimLeft(fullUrl.Path, "/")

				urldb := storage.GetUrldb()

				t.Log("Urldb: ", urldb)
				t.Log("request url: ", (*urldb)[shortUrl]+"/"+shortUrl)
				t.Log("body url the same: ", tt.body+shortUrl)

				//check request url and body url the same
				assert.Equal(t, (*urldb)[shortUrl]+"/"+shortUrl, tt.body+"/"+shortUrl)

				//go test -v ./...

			} else if tt.method == http.MethodGet {
				t.Log("=============GET===============")
				//get value of short URL from dburl:
				urldb := storage.GetUrldb()
				shortUrl, error := getShortUrl((*urldb), tt.body)

				t.Log("shortUrl: ", shortUrl)
				require.NotNil(t, error)

				//
				requestUrl := tt.request + "/" + shortUrl
				t.Log("requestUrl: ", requestUrl)

				request := httptest.NewRequest(http.MethodGet, requestUrl, nil)
				//create status recorder
				resRecord := httptest.NewRecorder()
				webhandl.GetUrl(resRecord, request)

				//get result
				res := resRecord.Result()
				//check answer code
				t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
				assert.Equal(t, tt.statusCode, res.StatusCode)

				//check Content type
				t.Log("Content-Type exp: ", tt.contentType, " act: ", res.Header.Get("Content-Type"))
				assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))

			}
		})

	}
}
