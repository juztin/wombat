package backends

import (
	"fmt"
)

type Backend interface{}

/* --------------------------------- Error ---------------------------------- */
type Status int

const (
	//StatusUnknown Status = 1000 + iota
	StatusDatastoreError = 1000 + iota
	StatusNotFound
	StatusNotModified
)

type Error interface {
	error
	Status() Status
	Critical() bool
}

type err struct {
	status Status
	msg    string
}

func (b *err) Error() string {
	return b.msg
}

func (b *err) Status() Status {
	return b.status
}

func (b *err) Critical() bool {
	return b.status <= StatusDatastoreError
}

func NewError(s Status, f string, params ...interface{}) Error {
	msg := fmt.Sprintf("[backend, %v] %s %v", s, f, params)
	return &err{s, msg}
}
