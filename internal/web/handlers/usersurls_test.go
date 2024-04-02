package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestUsersUrls(t *testing.T) {
	tests := []struct {
		name        string
		multipleURL string
		usersURLs   string
		numURLs     int
		statusCode  int
	}{
		{
			name:        "base test POTS",
			multipleURL: "http://localhost:8080/api/shorten/batch",
			usersURLs:   "http://localhost:8080/api/user/urls",
			numURLs:     20,
			statusCode:  http.StatusCreated,
		},
	}

	// init configApp
	app.InitLog()

	// init configApp
	configApp := config.DefaultConfig()

	// init config with defaults values

	short := service.NewService(storage.NewMemory())

	// init storage
	apiBatch := NewHandlerBatch(&configApp, short)

	// Get all users URLs.
	apiUsersURLs := NewHandlerAuthUser(&configApp, short)
	userID, err := uuid.NewV7()
	if err != nil {
		zap.S().Errorln("Error generate user uuid")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var insertURLS []entities.BatchRequest
			for i := 0; i < tt.numURLs; i++ {
				insertURLS = append(insertURLS, entities.BatchRequest{SessionID: strconv.Itoa(i), Origin: "http://yandex" + strconv.Itoa(i) + ".ru"})
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
			//check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)

			// unmarshal body
			var resp []entities.BatchResponse
			err = json.NewDecoder(res.Body).Decode(&resp)
			require.NoError(t, err)

			// check short URLS

			// add chi context
			rctx = chi.NewRouteContext()
			req = httptest.NewRequest(http.MethodGet, tt.usersURLs, nil)
			ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userID.String(), false))
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			req.Header.Add("Content-Type", "application/json")
			cookie = http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)

			// create status recorder
			resRecord = httptest.NewRecorder()
			apiUsersURLs.GetUserURLs(resRecord, req)
			// get result
			res = resRecord.Result()
			defer res.Body.Close()

			resAuth := []ResponseAuth{}

			err = json.NewDecoder(res.Body).Decode(&resAuth)
			require.NoError(t, err)

			// check answer code
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tt.numURLs, len(resAuth))
		})
	}
}
