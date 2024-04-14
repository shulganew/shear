package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAPIStat(t *testing.T) {
	tests := []struct {
		name       string
		trusted    bool
		shorts     int
		users      int
		err        error
		statusCode int
	}{
		{
			name:       "Get Stat",
			trusted:    true,
			shorts:     111,
			users:      222,
			err:        nil,
			statusCode: http.StatusOK,
		},
		{
			name:       "Get WithError",
			trusted:    true,
			shorts:     0,
			users:      0,
			err:        errors.New("Connection error"),
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Get WithError",
			trusted:    false,
			shorts:     0,
			users:      0,
			err:        errors.New("Connection error"),
			statusCode: http.StatusForbidden,
		},
	}
	app.InitLog()

	for _, tt := range tests {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// crete mock storege
		storeMock := mocks.NewMockStorageURL(ctrl)

		// init configApp
		configApp := config.DefaultConfig(false)

		// init storage
		handler := NewHandlerStat(&configApp, service.NewService(storeMock))

		_ = storeMock.EXPECT().
			GetNumShorts(gomock.Any()).
			AnyTimes().
			Return(tt.shorts, tt.err)

		_ = storeMock.EXPECT().
			GetNumUsers(gomock.Any()).
			AnyTimes().
			Return(tt.users, tt.err)

		t.Log(tt.name)

		// add chi context
		rctx := chi.NewRouteContext()

		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)

		ctxUser := context.WithValue(req.Context(), config.CtxAllow{}, tt.trusted)

		req = req.WithContext(context.WithValue(ctxUser, chi.RouteCtxKey, rctx))

		// create status recorder
		resRecord := httptest.NewRecorder()

		handler.GetStat(resRecord, req)

		// get result
		res := resRecord.Result()

		defer res.Body.Close()
		t.Log("StatusCode test: ", tt.statusCode, " server: ", res.StatusCode)
		assert.Equal(t, tt.statusCode, res.StatusCode)

		if tt.err == nil {
			var stat StatDTO
			err := json.NewDecoder(res.Body).Decode(&stat)
			assert.NoError(t, err)
			assert.Equal(t, tt.shorts, stat.Shorts)
			assert.Equal(t, tt.users, stat.Users)
		}
	}
}
