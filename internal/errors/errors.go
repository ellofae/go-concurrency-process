package errors

import "errors"

var (
	ErrNoRecordFound       = errors.New("no record has been found")
	ErrRecordAlreadyExists = errors.New("record already exists")
)
