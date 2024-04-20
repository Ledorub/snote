package response

import (
	"log"
	"net/http"
)

const (
	contentTypeJSON string = "application/json"
)

type JSONResponseWriter struct {
	logger  *log.Logger
	encoder jsonEncoder
}

func (writer *JSONResponseWriter) Write(w http.ResponseWriter, r *http.Request, status int, message any) {
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

func (writer *JSONResponseWriter) WriteBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errors := errorList{
		{"message": err.Error()},
	}
	writer.WriteError(w, r, http.StatusBadRequest, errors)
}

func (writer *JSONResponseWriter) WriteValidationError(w http.ResponseWriter, r *http.Request, errs []map[string]string) {
	errors := make(errorList, len(errs))
	for i, err := range errs {
		errors[i] = make(map[string]any, len(err))
		for k, v := range err {
			errors[i][k] = v
		}
	}
	writer.WriteError(w, r, http.StatusUnprocessableEntity, errors)
}

func NewJSONWriter(logger *log.Logger, encoder jsonEncoder) *JSONResponseWriter {
	return &JSONResponseWriter{logger: logger, encoder: encoder}
}

type jsonEncoder interface {
	Encode(data any) ([]byte, error)
}

type errorList = []map[string]any

func errorListToMsg(errors errorList) map[string]any {
	return map[string]any{
		"details": errors,
	}
}
