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
	app.InitLog()

	t.Run("db ping", func(t *testing.T) {
		sql, smock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}
		smock.ExpectPing()
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
		// check answer code
		t.Log("StatusCode test: ", http.StatusOK, " server: ", res.StatusCode)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
