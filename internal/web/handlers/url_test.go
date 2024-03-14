package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestURL(t *testing.T) {
	tests := []struct {
		name              string
		request           string
		body              string
		origin            string
		method            string
		contentType       string
		brief             string
		statusCode        int
		responseExist     bool
		responseIsDeleted bool
	}{
		{
			name:              "Set URL",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusCode:        307,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
		},

		{
			name:              "Get URL",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusCode:        410,
			brief:             "asdfghjk",
			responseExist:     true,
			responseIsDeleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// crete mock storege
			storeMock := mocks.NewMockStorageURL(ctrl)

			// init configApp
			configApp := &config.Config{}

			// init config with difauls values
			configApp.Address = config.DefaultHost
			configApp.Response = config.DefaultHost

			// init storage
			handler := NewHandlerGetURL(configApp, service.NewService(storeMock))

			userID, err := uuid.NewV7()
			if err != nil {
				zap.S().Errorln("Error generate user uuid")
			}

			_ = storeMock.EXPECT().
				GetOrigin(gomock.Any(), tt.brief).
				Times(1).
				Return(tt.origin, tt.responseExist, tt.responseIsDeleted)

			//
			requestURL, _ := url.JoinPath(tt.request, tt.brief)
			t.Log("requestUrl: ", requestURL)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.brief)

			// use context for chi router - add id
			req := httptest.NewRequest(http.MethodGet, requestURL, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)

			// create status recorder
			resRecord := httptest.NewRecorder()
			handler.GetURL(resRecord, req)

			// get result
			res := resRecord.Result()
			defer res.Body.Close()
			// check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)

			// check Content type
			t.Log("Content-Type exp: ", tt.contentType, " act: ", res.Header.Get("Content-Type"))
			assert.Equal(t, tt.contentType, res.Header.Get("Content-Type"))
		})
	}
}
