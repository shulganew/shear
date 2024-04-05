package entities

import (
	"database/sql"
)

// Base stract for working with storage and api. It contains main fields:
//
// • ID in database;
//
// • UUID of short object
//
// • Brief URL - shor value of user's URL
//
// • Origin - user's original long URL.
type Short struct {
	Brief     string         `json:"short_url"`    // short URL (cache)
	Origin    string         `json:"original_url"` // Long full URL
	SessionID string         `json:"-"`            // for Batch require use: Unique Session ID for each request in URL Batch
	UserID    sql.NullString `json:"-"`            // unique user id string
	IsDeleted bool           `json:"-"`            // mark deleted URL by user
	ID        int            `json:"uuid"`
}

// Constructor Short.
func NewShort(ID int, UUID string, brief string, origin string, sessionID string) *Short {
	nullUUID := sql.NullString{String: UUID, Valid: true}
	if UUID == "" {
		nullUUID.Valid = false
	}
	return &Short{ID: ID, IsDeleted: false, UserID: nullUUID, Brief: brief, Origin: origin, SessionID: sessionID}
}
