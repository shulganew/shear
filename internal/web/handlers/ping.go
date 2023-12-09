package handlers

import (
	"database/sql"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
)

// hadler for testing db connection

type Ping struct {
	conf *config.App
	db   *sql.DB
}

// Test DB connection
func (b *Ping) Ping(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Content-Type", "text/plain")

	if err := b.db.PingContext(req.Context()); err == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("<h1>Connected to Data Base!</h1>"))
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}

}

func NewDB(configApp *config.App) *Ping {

	return &Ping{conf: configApp, db: configApp.DB}
}
