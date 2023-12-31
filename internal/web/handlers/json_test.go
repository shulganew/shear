package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
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
	configApp.Address = config.DefaultHost
	configApp.Response = config.DefaultHost

	stor := service.StorageURL(storage.NewMemory())
	//init storage
	apiHand := NewHandlerAPI(configApp, &stor)
	serviceURL := apiHand.GetService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Log("=============POTS===============")
			t.Log("tt.request=", tt.body)

			//add chi context
			rctx := chi.NewRouteContext()

			req := httptest.NewRequest(http.MethodPost, tt.requestURL, strings.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req.Header.Add("Content-Type", "application/json")
			//create status recorder
			resRecord := httptest.NewRecorder()

			apiHand.GetBrief(resRecord, req)

			//get result
			res := resRecord.Result()
			defer res.Body.Close()
			//check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)

			//Unmarshal body
			var response Resonse

			err := json.NewDecoder(res.Body).Decode(&response)
			require.NoError(t, err)

			//responseURL = hostname+brief
			responseURL, err := url.Parse(response.Brief)
			require.NoError(t, err)
			t.Log(responseURL)
			brief := strings.TrimLeft(responseURL.Path, "/")

			originDB, exist := serviceURL.GetOrigin(req.Context(), brief)
			require.True(t, exist)

			t.Log("brief url: ", originDB)

			assert.Equal(t, originDB, tt.link)

		})
	}
}
