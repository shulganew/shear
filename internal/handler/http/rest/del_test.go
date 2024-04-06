package rest

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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	"github.com/shulganew/shear.git/internal/handler/http/middlewares"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type UserDelSet struct {
	userID    string
	delShorts []string
}

func TestDelBulk(t *testing.T) {
	tests := []struct {
		name        string
		multipleURL string
		delURL      string
		reqestShort string
		numUsers    int
		numURLs     int
	}{
		{
			name:        "Create del pool",
			multipleURL: "http://localhost:8080/api/shorten/batch",
			delURL:      "http://localhost:8080/api/user/urls",
			reqestShort: "http://localhost:8080/",
			numUsers:    3,
			numURLs:     20,
		},
	}

	// init configApp
	app.InitLog()
	// init configApp
	configApp := config.DefaultConfig(false)

	stor := storage.NewMemory()
	short := service.NewService(stor)
	// init storage
	apiBatch := NewHandlerBatch(&configApp, short)
	handGet := NewHandlerGetURL(&configApp, short)
	delCh := make(chan service.DelBatch, 100)
	defer close(delCh)
	del := service.NewDelete(delCh, &configApp)
	handDel := NewHandlerDelShorts(del)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Create UserDelSet of the new URLS
			userSet := make([]UserDelSet, tt.numUsers)
			for i := 0; i < tt.numUsers; i++ {
				userID, err := uuid.NewV7()
				userSet[i] = UserDelSet{userID: userID.String()}
				if err != nil {
					zap.S().Errorln("Error generate user uuid")
				}

				var insertURLS []entities.BatchRequest
				for j := 0; j < tt.numURLs; j++ {
					insertURLS = append(insertURLS, entities.BatchRequest{SessionID: strconv.Itoa(j), Origin: "http://yandex" + strconv.Itoa(j) + ".ru"})
				}

				body, err := json.Marshal(&insertURLS)
				require.NoError(t, err)

				// add chi context
				rctx := chi.NewRouteContext()
				t.Log("URL: ", tt.multipleURL)
				req := httptest.NewRequest(http.MethodPost, tt.multipleURL, bytes.NewReader(body))
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
				req.Header.Add("Content-Type", "application/json")
				cookie := http.Cookie{Name: "user_id", Value: userSet[i].userID}
				req.AddCookie(&cookie)

				// create status recorder
				resRecord := httptest.NewRecorder()
				apiBatch.BatchAdd(resRecord, req)

				// 2. Reas body for result and check server answer
				res := resRecord.Result()
				defer res.Body.Close()
				// check answer code
				t.Log("StatusCode test: ", http.StatusCreated, " server: ", res.StatusCode)
				assert.Equal(t, http.StatusCreated, res.StatusCode)

				// 3. Save Users shorts to UserDelSet
				var resp []entities.BatchResponse
				err = json.NewDecoder(res.Body).Decode(&resp)
				require.NoError(t, err)

				for _, short := range resp {
					URL, err := url.Parse(short.Answer)
					require.NoError(t, err)
					path := URL.Path
					short := strings.TrimPrefix(path, "/")
					userSet[i].delShorts = append(userSet[i].delShorts, short)
				}
			}

			// 4.Make buch apdate
			for i := 0; i < tt.numUsers; i++ {
				body, err := json.Marshal(&userSet[i].delShorts)

				t.Log("BODY: ", string(body))
				require.NoError(t, err)
				// add chi context
				rctx := chi.NewRouteContext()
				t.Log("URL: ", tt.delURL)

				req := httptest.NewRequest(http.MethodDelete, tt.delURL, bytes.NewReader(body))
				ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userSet[i].userID, false))
				req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

				req.Header.Add("Content-Type", "application/json")
				cookie := http.Cookie{Name: "user_id", Value: userSet[i].userID}
				req.AddCookie(&cookie)
				// create status recorder
				resRecord := httptest.NewRecorder()
				handDel.DelUserURLs(resRecord, req)
				res := resRecord.Result()
				defer res.Body.Close()
				// check answer code
				t.Log("StatusCode test: ", http.StatusAccepted, " server: ", res.StatusCode)
				assert.Equal(t, http.StatusAccepted, res.StatusCode)
				// 5. Wait, then check if ULR field change to is_delete == true
			}
			time.Sleep(time.Second)
			// Check all created del request in main feth GET API (GetURL)
			for i := 0; i < tt.numUsers; i++ {
				for j := 0; j < tt.numURLs; j++ {
					//add chi context
					rctx := chi.NewRouteContext()
					rctx.URLParams.Add("id", userSet[i].delShorts[j])
					req := httptest.NewRequest(http.MethodGet, tt.reqestShort, nil)
					ctx := context.WithValue(req.Context(), config.CtxPassKey{}, configApp.GetPass())
					req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
					req.Header.Add("Content-Type", "plain/text")
					cookie := http.Cookie{Name: "user_id", Value: userSet[i].userID}
					req.AddCookie(&cookie)

					// create status recorder
					resRecord := httptest.NewRecorder()
					// Add midleware manualy
					getHandler := middlewares.Auth(http.HandlerFunc(handGet.GetURL))
					getHandler.ServeHTTP(resRecord, req)
					// get result
					res := resRecord.Result()
					defer res.Body.Close()
					// check answer code
					t.Log("StatusCode test: ", http.StatusTemporaryRedirect, " server: ", res.StatusCode)
					assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
				}
			}
		})
	}
}
