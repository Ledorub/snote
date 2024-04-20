package common

import (
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
	WriteValidationError(http.ResponseWriter, *http.Request, ValidationErrors)
}

type Validator interface {
	CheckIsValid() bool
	AddNonFieldError(string)
	AddFieldError(string, string)
	CheckField(string, bool, string)
	GetNonFieldErrors() []string
	GetFieldErrors() map[string]string
}

type ValidatorFactory = func() Validator

type ValidationErrors = []map[string]string
