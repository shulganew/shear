package storage

import (
	"fmt"

	"github.com/shulganew/shear.git/internal/service"
)

// Error use when ID unknown
type ErrDuplicatedURL struct {
	Label  string
	Brief  string
	Origin string
	Err    error
}

func (ed *ErrDuplicatedURL) Error() string {
	return fmt.Sprintf("%s. Existed URL: [%s]. Error: [%v]", ed.Label, ed.Brief, ed.Err)
}

// Unwrap()
func (ed *ErrDuplicatedURL) Unwrap() error {
	return ed.Err
}

func NewErrDuplicatedURL(brief string, origin string, err error) *ErrDuplicatedURL {
	return &ErrDuplicatedURL{Label: "URL alredy existed. ", Origin: origin, Brief: brief, Err: err}
}

// Error use when ID exit, creates object short
type ErrDuplicatedShort struct {
	Label string
	Short service.Short
	Err   error
}

func (ed *ErrDuplicatedShort) Error() string {
	return fmt.Sprintf("%s. Existed URL: [%s]. Error: [%v]", ed.Label, ed.Short.Brief, ed.Err)
}

// Unwrap()
func (ed *ErrDuplicatedShort) Unwrap() error {
	return ed.Err
}

func NewErrDuplicatedShort(sessionID string, brief string, origin string, err error) *ErrDuplicatedShort {
	return &ErrDuplicatedShort{Label: "URL alredy existed. ", Short: service.Short{SessionID: sessionID, Brief: brief, Origin: origin}, Err: err}
}
