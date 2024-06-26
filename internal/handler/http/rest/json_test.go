package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAPI(t *testing.T) {
	tests := []struct {
		name        string
		requestURL  string
		body        string
		link        string
		statusCode  int
		isCookie    bool
		isBrokenURL bool
		//want
	}{
		{
			name:        "base test POTS",
			requestURL:  "http://localhost:8080/api/shorten",
			body:        `{"url": "https://practicum.yandex.ru"}`,
			link:        "https://practicum.yandex.ru",
			statusCode:  http.StatusCreated,
			isCookie:    true,
			isBrokenURL: false,
		},

		{
			name:        "base test Broken JSON",
			requestURL:  "http://localhost:8080/api/shorten",
			body:        `{"ul": "https://practicum.yandex.ru"}`,
			link:        "https://practicum.yandex.ru",
			statusCode:  http.StatusInternalServerError,
			isCookie:    true,
			isBrokenURL: true,
		},
	}
	app.InitLog()
	// init configApp
	configApp := config.DefaultConfig(false)

	short := service.NewService(storage.NewMemory())
	// init storage
	apiHand := NewHandlerAPI(&configApp, short)

	userID, err := uuid.NewV7()
	if err != nil {
		zap.S().Errorln("Error generate user uuid")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("tt.request=", tt.body)

			// add chi context
			rctx := chi.NewRouteContext()

			req := httptest.NewRequest(http.MethodPost, tt.requestURL, strings.NewReader(tt.body))
			ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userID.String(), !tt.isCookie))
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			req.Header.Add("Content-Type", "application/json")
			if tt.isCookie {
				cookie := http.Cookie{Name: "user_id", Value: userID.String()}
				req.AddCookie(&cookie)
			}

			// create status recorder
			resRecord := httptest.NewRecorder()

			apiHand.GetBrief(resRecord, req)

			// get result
			res := resRecord.Result()
			defer res.Body.Close()

			if !tt.isCookie {
				// check answer code
				t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
				assert.Equal(t, tt.statusCode, res.StatusCode)
				return
			}

			// check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)

			if !tt.isBrokenURL {
				// unmarshal body
				var response Response
				err := json.NewDecoder(res.Body).Decode(&response)

				require.NoError(t, err)

				// responseURL = hostname+brief
				responseURL, err := url.Parse(response.Brief)
				require.NoError(t, err)
				t.Log(responseURL)
				brief := strings.TrimLeft(responseURL.Path, "/")

				originDB, exist, _ := short.GetOrigin(req.Context(), brief)
				require.True(t, exist)

				t.Log("brief url: ", originDB)
				assert.Equal(t, originDB, tt.link)
			}
		})
	}
}
