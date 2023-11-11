package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/app/config"
	utils "github.com/shulganew/shear.git/internal/core"
	webhandl "github.com/shulganew/shear.git/internal/handlers"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	// init configApp
	configApp := config.GetConfig()
	// init config
	configApp.StartAddress = config.DefaultHost
	configApp.ResultAddress = config.DefaultHost

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

				//responseUrl = hostname+shortUrl
				responseUrl, err := url.Parse(body)
				require.NoError(t, err)

				t.Log("full url: ", responseUrl.Path)

				shortUrl := strings.TrimLeft(responseUrl.Path, "/")
				urldb := *storage.GetUrldb()
				longUrlDb := urldb[shortUrl]
				t.Log("shortUrl url: ", shortUrl)
				responseUrlDb, err := url.JoinPath(longUrlDb.String(), shortUrl)
				require.NoError(t, err)

				t.Log("Urldb: ", urldb)
				t.Log("ressponseUrl from db: ", responseUrlDb)
				bodyUrl, _ := url.JoinPath(tt.body, shortUrl)
				t.Log("body url the same: ", bodyUrl)

				//check request url and body url the same
				assert.Equal(t, responseUrlDb, bodyUrl)

				//go test -v ./...

			} else if tt.method == http.MethodGet {
				t.Log("=============GET===============")
				//get value of short URL from urldb:
				urldb := storage.GetUrldb()

				shortUrl, error := utils.GetShortUrl(urldb, tt.body)

				t.Log("shortUrl: ", shortUrl)
				require.NotNil(t, error)

				//
				requestUrl, _ := url.JoinPath(tt.request, shortUrl)
				t.Log("requestUrl: ", requestUrl)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", shortUrl)

				//use context for chi router - add id
				request := httptest.NewRequest(http.MethodGet, requestUrl, nil)
				request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

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
