package note

import (
	"github.com/ledorub/snote-api/api/common"
	"log"
	"net/http"
	"time"
)

type API struct {
	logger         *log.Logger
	requestReader  common.RequestReader
	responseWriter common.ResponseWriter
}

func New(logger *log.Logger, requestReader common.RequestReader, responseWriter common.ResponseWriter) *API {
	return &API{logger: logger, requestReader: requestReader, responseWriter: responseWriter}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	type noteCreate struct {
		Content           string        `json:"content"`
		ExpiresAt         time.Time     `json:"expiresAt"`
		ExpiresAtTimezone string        `json:"expiresAtTimezone"`
		ExpiresIn         time.Duration `json:"expiresIn"`
		KeyHash           string        `json:"keyHash"`
	}
	note := noteCreate{}
	if err := api.requestReader.Read(r.Body, &note); err != nil {
		api.responseWriter.WriteBadRequest(w, r, err)
		return
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
