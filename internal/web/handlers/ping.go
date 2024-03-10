package handlers

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"
)

// Handler for testing db connection.
//
//	Get "/ping"
type Ping struct {
	db *sql.DB
}

// Service constructor.
func NewDB(db *sql.DB) *Ping {

	return &Ping{db: db}
}

// Test DB connection.
// @Summary      Test database
// @Description  Ping service for database connection check
// @Tags         api
// @Success      200 "Available"
// @Failure      500 "Handling error"
// @Router       /ping [get]
func (b *Ping) Ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain")
	if err := b.db.PingContext(req.Context()); err == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("<h1>Connected to Data Base!</h1>"))
	} else {
		zap.S().Errorln(err)
		res.WriteHeader(http.StatusInternalServerError)
	}
}
