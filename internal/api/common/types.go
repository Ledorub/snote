package common

import (
	"context"
	"github.com/ledorub/snote-api/internal"
	"github.com/ledorub/snote-api/internal/validator"
	"io"
	"net/http"
)

type JSONEncDec interface {
	Encode(data map[string]any) ([]byte, error)
	Decode(data io.Reader, dst any) error
}

type RequestReader interface {
	Read(body io.Reader, dst any) error
}

type ResponseWriter interface {
	Write(http.ResponseWriter, *http.Request, int, any)
	WriteError(http.ResponseWriter, *http.Request, int, []map[string]any)
	WriteServerError(http.ResponseWriter, *http.Request, error)
	WriteNotFound(http.ResponseWriter, *http.Request)
	WriteBadRequest(http.ResponseWriter, *http.Request, error)
	WriteValidationError(http.ResponseWriter, *http.Request, []error)
}

type Validator interface {
	Check(bool, string)
	CheckIsValid() bool
	GetErrors() []validator.ValidationError
}

type ValidatorFactory = func() Validator

type NoteService interface {
	CreateNote(ctx context.Context, note *internal.Note) (*internal.Note, error)
}
