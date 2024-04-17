package common

import (
	"io"
	"net/http"
)

type JSONEncDec interface {
	Encode(data map[string]any) ([]byte, error)
	Decode(data io.Reader, dst any) error
}

type ResponseWriter interface {
	Write(w http.ResponseWriter, r *http.Request, status int, message map[string]any)
	WriteError(w http.ResponseWriter, r *http.Request, status int, errors []map[string]any)
	WriteServerError(w http.ResponseWriter, r *http.Request, err error)
	WriteNotFound(w http.ResponseWriter, r *http.Request)
}
