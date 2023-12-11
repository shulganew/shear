package storage

import "fmt"

// custom error
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
