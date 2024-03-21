package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const DefaultHost = "localhost:8080"

func TestBatch(t *testing.T) {
	tests := []struct {
		name        string
		multipleURL string
		reqestShort string
		numURLs     int
		statusCode  int
		brokenURL   bool
	}{
		{
			name:        "base test POTS",
			multipleURL: "http://localhost:8080/api/shorten/batch",
			reqestShort: "http://localhost:8080/",
			numURLs:     20,
			statusCode:  http.StatusCreated,
			brokenURL:   false,
		},
		{
			name:        "Internal error - parse broken URL",
			multipleURL: "http://localhost:8080/api/shorten/batch",
			reqestShort: "http://localhost:8080/",
			numURLs:     20,
			statusCode:  http.StatusInternalServerError,
			brokenURL:   true,
		},
	}

	// init configApp
	app.InitLog()
	// init configApp
	configApp := &config.Config{}
	// init config with difauls values
	configApp.Address = DefaultHost
	configApp.Response = DefaultHost
	configApp.IsDB = false
	configApp.IsBackup = false
	short := service.NewService(storage.NewMemory())
	// init storage
	apiBatch := NewHandlerBatch(configApp, short)
	webHand := NewHandlerGetURL(configApp, short)

	userID, err := uuid.NewV7()
	if err != nil {
		zap.S().Errorln("Error generate user uuid")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var insertURLS []entities.BatchRequest
			if !tt.brokenURL {
				for i := 0; i < tt.numURLs; i++ {
					insertURLS = append(insertURLS, entities.BatchRequest{SessionID: strconv.Itoa(i), Origin: "http://yandex" + strconv.Itoa(i) + ".ru"})
				}

			} else {
				for i := 0; i < tt.numURLs; i++ {
					insertURLS = append(insertURLS, entities.BatchRequest{SessionID: strconv.Itoa(i), Origin: "yandex" + strconv.Itoa(i) + "ru"})
				}
			}
			body, err := json.Marshal(&insertURLS)
			require.NoError(t, err)

			// add chi context
			rctx := chi.NewRouteContext()
			t.Log("URL: ", tt.multipleURL)
			req := httptest.NewRequest(http.MethodPost, tt.multipleURL, bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req.Header.Add("Content-Type", "application/json")
			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)

			// create status recorder
			resRecord := httptest.NewRecorder()
			apiBatch.BatchSet(resRecord, req)

			// get result
			res := resRecord.Result()
			defer res.Body.Close()
			// check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)
			if tt.brokenURL {
				return
			}
			// unmarshal body
			var resp []entities.BatchResponse

			err = json.NewDecoder(res.Body).Decode(&resp)
			require.NoError(t, err)

			// check short URLS
			for _, short := range resp {

				// add chi context
				rctx = chi.NewRouteContext()
				URL, err := url.Parse(short.Answer)
				require.NoError(t, err)
				id := URL.Path

				rctx.URLParams.Add("id", strings.TrimPrefix(id, "/"))
				req = httptest.NewRequest(http.MethodGet, short.Answer, nil)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				req.Header.Add("Content-Type", "plain/text")
				cookie = http.Cookie{Name: "user_id", Value: userID.String()}
				req.AddCookie(&cookie)

				// create status recorder
				resRecord = httptest.NewRecorder()
				webHand.GetURL(resRecord, req)

				// get result
				res := resRecord.Result()
				defer res.Body.Close()
				// check answer code
				assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
			}
		})
	}
}
