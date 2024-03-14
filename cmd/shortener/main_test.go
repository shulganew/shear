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
	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/web/handlers"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
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
	configApp := config.InitConfig()

	// init config with difauls values
	configApp.Address = config.DefaultHost
	configApp.Response = config.DefaultHost

	short := service.NewService(storage.NewMemory())

	//init storage
	handler := handlers.NewHandlerGetURL(configApp, short)

	userID, err := uuid.NewV7()
	if err != nil {
		zap.S().Errorln("Error generate user uuid")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.method == http.MethodPost {
				t.Log("=============POTS===============")
				t.Log("tt.request=", tt.request)
				t.Log("strings.NewReader(tt.body)=", tt.body)
				req := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
				cookie := http.Cookie{Name: "user_id", Value: userID.String()}
				req.AddCookie(&cookie)
				//create status recorder
				resRecord := httptest.NewRecorder()

				handler.SetURL(resRecord, req)

				//get result
				res := resRecord.Result()

				//check answer code
				t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
				assert.Equal(t, tt.statusCode, res.StatusCode)

				//check Content type
				t.Log("Content-Type exp: ", tt.contentType, " act: ", res.Header.Get("Content-Type"))
				assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))

				//check body content

				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				body := string(resBody)
				t.Log("Body: ", body)
				defer res.Body.Close()

				//responseURL = hostname+brief
				responseURL, err := url.Parse(body)
				require.NoError(t, err)

				t.Log("full url: ", responseURL.Path)

				brief := strings.TrimLeft(responseURL.Path, "/")

				originDB, exist, _ := short.GetOrigin(req.Context(), brief)
				require.True(t, exist)

				t.Log("brief url: ", brief)
				responseURLDb, err := url.JoinPath(originDB, brief)
				require.NoError(t, err)

				t.Log("ressponseUrl from db: ", responseURLDb)
				bodyURL, _ := url.JoinPath(tt.body, brief)
				t.Log("body url the same: ", bodyURL)

				//check request url and body url the same
				assert.Equal(t, responseURLDb, bodyURL)

			} else if tt.method == http.MethodGet {
				t.Log("=============GET===============")

				//get brief from storage
				brief, err, _ := short.GetBrief(context.Background(), tt.body)

				t.Log("brief: ", brief)
				require.NotNil(t, err)

				//
				requestURL, _ := url.JoinPath(tt.request, brief)
				t.Log("requestUrl: ", requestURL)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", brief)

				//use context for chi router - add id
				req := httptest.NewRequest(http.MethodGet, requestURL, nil)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				cookie := http.Cookie{Name: "user_id", Value: userID.String()}
				req.AddCookie(&cookie)

				//create status recorder
				resRecord := httptest.NewRecorder()
				handler.GetURL(resRecord, req)

				//get result
				res := resRecord.Result()
				defer res.Body.Close()
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
