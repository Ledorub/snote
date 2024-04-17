package response

import (
	"log"
	"net/http"
)

const (
	contentTypeJSON string = "application/json"
)

type errorList = []map[string]any

type JSONEncoder interface {
	Encode(data map[string]any) ([]byte, error)
}

type JSONResponseWriter struct {
	logger  *log.Logger
	encoder JSONEncoder
}

func (writer *JSONResponseWriter) Write(w http.ResponseWriter, r *http.Request, status int, message map[string]any) {
	encoded, err := writer.encoder.Encode(message)
	if err != nil {
		writer.logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)
	w.Write(encoded)
}

func (writer *JSONResponseWriter) WriteError(w http.ResponseWriter, r *http.Request, status int, errors errorList) {
	message := errorListToMsg(errors)
	writer.Write(w, r, status, message)
}

func (writer *JSONResponseWriter) WriteServerError(w http.ResponseWriter, r *http.Request, err error) {
	writer.logger.Print(err)
	errors := errorList{
		{"message": err.Error()},
	}
	writer.WriteError(w, r, http.StatusInternalServerError, errors)
}

func (writer *JSONResponseWriter) WriteNotFound(w http.ResponseWriter, r *http.Request) {
	errors := errorList{
		{"message": "Not Found"},
	}
	writer.WriteError(w, r, http.StatusNotFound, errors)
}

func NewJSONWriter(logger *log.Logger, encoder JSONEncoder) *JSONResponseWriter {
	return &JSONResponseWriter{logger: logger, encoder: encoder}
}

func errorListToMsg(errors errorList) map[string]any {
	return map[string]any{
		"details": errors,
	}
}
