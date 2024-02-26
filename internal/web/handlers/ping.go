package handlers

import (
	"database/sql"
	"net/http"
)

// Handler for testing db connection.
//
//	Get "/ping"
type Ping struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *Ping {

	return &Ping{db: db}
}

// Test DB connection.
func (b *Ping) Ping(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain")
	if err := b.db.PingContext(req.Context()); err == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("<h1>Connected to Data Base!</h1>"))
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
