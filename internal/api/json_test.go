package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_api(t *testing.T) {
	tests := []struct {
		name       string
		requestURL string
		body       string
		link       string
		statusCode int
		//want
	}{
		{
			name:       "base test POTS",
			requestURL: "http://localhost:8080/api/shorten",
			body:       `{"url": "https://practicum.yandex.ru"}`,
			link:       "https://practicum.yandex.ru",
			statusCode: http.StatusCreated,
		},

		{
			name:       "base test POTS",
			requestURL: "http://localhost:8080/api/shorten",
			body:       `{"url": "https://practicum.yandex.ru"}`,
			link:       "https://practicum.yandex.ru",
			statusCode: http.StatusCreated,
		},
	}
	// init configApp
	configApp := config.InitConfig()
	// init config with difauls values
	configApp.StartAddress = config.DefaultHost
	configApp.ResultAddress = config.DefaultHost
	configApp.Storage = &storage.MapStorage{StoreURLs: make(map[string]url.URL)}

	//init storage
	apiHand := NewHandler(configApp)
	serviceURL := apiHand.GetService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Log("=============POTS===============")
			t.Log("tt.request=", tt.body)

			//add chi context
			rctx := chi.NewRouteContext()

			req := httptest.NewRequest(http.MethodPost, tt.requestURL, strings.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req.Header.Add("Content-Type", CONTENT_TYPE_JSON)
			//create status recorder
			resRecord := httptest.NewRecorder()

			apiHand.SetAPI(resRecord, req)

			//get result
			res := resRecord.Result()

			//check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)

			//Unmarshal body
			var response ResonseJSON

			err := json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			//responseURL = hostname+shortUrl
			responseURL, err := url.Parse(response.ShortURL)
			require.NoError(t, err)
			t.Log(responseURL)
			shortURL := strings.TrimLeft(responseURL.Path, "/")

			longURLDb, exist := serviceURL.GetLongURL(shortURL)
			require.True(t, exist)

			t.Log("shortUrl url: ", longURLDb)

			assert.Equal(t, longURLDb.String(), tt.link)

		})
	}
}
