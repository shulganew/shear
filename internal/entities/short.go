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
	ID        int            `json:"uuid"`
	UUID      sql.NullString `json:"-"`            // unique user id string
	Brief     string         `json:"short_url"`    // short URL (cache)
	Origin    string         `json:"original_url"` // Long full URL
	SessionID string         `json:"-"`            // for Batch reques use: Unique Session ID for each request in URL Batch
	IsDeleted bool           `json:"-"`            // mark deleted URL by user
}

// Constructor Short.
func NewShort(ID int, UUID string, brief string, origin string, sessionID string) *Short {
	nullUUID := sql.NullString{String: UUID, Valid: true}
	if UUID == "" {
		nullUUID.Valid = false
	}
	return &Short{ID, nullUUID, brief, origin, sessionID, false}
}
