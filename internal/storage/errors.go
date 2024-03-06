package storage

import (
	"fmt"

	"github.com/shulganew/shear.git/internal/entities"
)

// Error use when ID unknown.
type ErrDuplicatedURL struct {
	Err    error
	Label  string
	Brief  string
	Origin string
}

// Error constructor.
func NewErrDuplicatedURL(brief string, origin string, err error) *ErrDuplicatedURL {
	return &ErrDuplicatedURL{Label: "URL already existed. ", Origin: origin, Brief: brief, Err: err}
}

// Return Error string description.
func (ed *ErrDuplicatedURL) Error() string {
	return fmt.Sprintf("%s. Existed URL: [%s]. Error: [%v]", ed.Label, ed.Brief, ed.Err)
}

// Unwrap error.
func (ed *ErrDuplicatedURL) Unwrap() error {
	return ed.Err
}

// Error use when ID exit, creates object short.
type ErrDuplicatedShort struct {
	Err   error
	Label string
	Short entities.Short
}

// Error constructor.
func NewErrDuplicatedShort(sessionID string, brief string, origin string, err error) *ErrDuplicatedShort {
	return &ErrDuplicatedShort{Label: "URL alredy existed. ", Short: entities.Short{SessionID: sessionID, Brief: brief, Origin: origin}, Err: err}
}

// Return Error string description.
func (ed *ErrDuplicatedShort) Error() string {
	return fmt.Sprintf("%s. Existed URL: [%s]. Error: [%v]", ed.Label, ed.Short.Brief, ed.Err)
}

// Unwrap error.
func (ed *ErrDuplicatedShort) Unwrap() error {
	return ed.Err
}
