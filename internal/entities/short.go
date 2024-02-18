package entities

import (
	"database/sql"
)

// base stract for working with storage
type Short struct {
	ID int `json:"uuid"`

	//unique user id string
	UUID sql.NullString `json:"-"`

	//mark deleted URL by user
	IsDeleted bool `json:"-"`

	//short URL (cache)
	Brief string `json:"short_url"`

	//Long full URL
	Origin string `json:"original_url"`

	//For Batch reques use: Unique Session ID for each request in URL Batch
	SessionID string `json:"-"`
}

func NewShort(ID int, UUID string, brief string, origin string, sessionID string) *Short {
	nullUUID := sql.NullString{String: UUID, Valid: true}
	if UUID == "" {
		nullUUID.Valid = false
	}

	return &Short{ID, nullUUID, false, brief, origin, sessionID}
}
