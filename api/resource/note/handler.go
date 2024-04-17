package note

import (
	"github.com/ledorub/snote-api/api/common"
	"log"
	"net/http"
)

type API struct {
	logger         *log.Logger
	responseWriter common.ResponseWriter
}

func New(logger *log.Logger, responseWriter common.ResponseWriter) *API {
	return &API{logger: logger, responseWriter: responseWriter}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	note := map[string]any{
		"id":      2,
		"content": "encrypted content",
	}
	api.responseWriter.Write(w, r, http.StatusCreated, note)
}

func (api *API) Read(w http.ResponseWriter, r *http.Request) {
	note := map[string]any{
		"id":      1,
		"content": "encrypted content",
	}
	api.responseWriter.Write(w, r, http.StatusOK, note)
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	api.responseWriter.Write(w, r, http.StatusNoContent, nil)
}
