package handlers

import (
	"database/sql"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
)

// hadler for testing db connection

type Base struct {
	conf *config.Shear
	db   *sql.DB
}

// Test DB connection
func (b *Base) Ping(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Content-Type", "text/plain")

	if err := b.db.Ping(); err == nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}

	res.Write([]byte("<h1>Connected to Data Base!</h1>"))

}

func NewDB(configApp *config.Shear) *Base {

	return &Base{conf: configApp, db: configApp.DB}
}
