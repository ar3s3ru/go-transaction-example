package app

import (
	"fmt"
	"net/http"
)

type Error interface {
	error

	Unwrap() error
	StatusCode() int
}

type AlreadyExistsError struct{ Inner error }

func (e AlreadyExistsError) Error() string   { return fmt.Sprintf("app: already exists (%s)", e.Inner) }
func (e AlreadyExistsError) Unwrap() error   { return e.Inner }
func (e AlreadyExistsError) StatusCode() int { return http.StatusConflict }

type NotFoundError struct{ Inner error }

func (e NotFoundError) Error() string   { return fmt.Sprintf("app: not found (%s)", e.Inner) }
func (e NotFoundError) Unwrap() error   { return e.Inner }
func (e NotFoundError) StatusCode() int { return http.StatusNotFound }
