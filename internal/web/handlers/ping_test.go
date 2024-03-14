package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestAPIPing(t *testing.T) {
	tests := []struct {
		name string
		ping bool
		//want
	}{
		{
			name: "DB ping test available",
			ping: true,
		},

		{
			name: "DB ping test not available",
			ping: false,
		},
	}
	app.InitLog()

	for _, tt := range tests {
		t.Log(tt.name)

		sql, smock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}
		if tt.ping {
			smock.ExpectPing()
		}
		pingSrv := NewDB(sql)

		// add chi context
		rctx := chi.NewRouteContext()

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// create status recorder
		resRecord := httptest.NewRecorder()

		pingSrv.Ping(resRecord, req)

		// get result
		res := resRecord.Result()
		defer res.Body.Close()
		if tt.ping {
			// check answer code
			t.Log("StatusCode test: ", http.StatusOK, " server: ", res.StatusCode)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		} else {
			t.Log("StatusCode test: ", http.StatusOK, " server: ", res.StatusCode)
			assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		}
	}
}
