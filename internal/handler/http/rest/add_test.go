package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAdd(t *testing.T) {
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
		err               error
	}{
		{
			name:              "Add URL",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusCode:        http.StatusCreated,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
			err:               nil,
		},
		{
			name:              "Add URL DB error",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusCode:        http.StatusInternalServerError,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
			err:               errors.New("Database error"),
		},

		{
			name:              "Add URL Duplication error",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusCode:        http.StatusConflict,
			brief:             "dupli234",
			responseExist:     true,
			responseIsDeleted: false,
			err:               &service.ErrDuplicatedURL{Err: errors.New("Database duplicated"), Label: "Duplicated", Brief: "dupli234", Origin: "http://localhost:8080"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// crete mock storege
			storeMock := mocks.NewMockStorageURL(ctrl)

			// init configApp
			configApp := config.DefaultConfig(false)

			// init storage
			handler := NewHandlerGetURL(&configApp, service.NewService(storeMock))

			userID, err := uuid.NewV7()
			if err != nil {
				zap.S().Errorln("Error generate user uuid")
			}

			_ = storeMock.EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(tt.err)

			rctx := chi.NewRouteContext()

			// use context for chi router - add id
			req := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			ctx := context.WithValue(req.Context(), config.CtxConfig{}, config.NewCtxConfig(userID.String(), false))
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			req.AddCookie(&cookie)

			// create status recorder
			resRecord := httptest.NewRecorder()
			handler.AddURL(resRecord, req)

			// get result
			res := resRecord.Result()
			defer res.Body.Close()

			// check answer code
			t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
			assert.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}
