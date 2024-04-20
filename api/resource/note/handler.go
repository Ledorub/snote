package note

import (
	"github.com/ledorub/snote-api/api/common"
	"log"
	"net/http"
	"time"
)

type API struct {
	logger           *log.Logger
	requestReader    common.RequestReader
	responseWriter   common.ResponseWriter
	validatorFactory common.ValidatorFactory
}

func New(logger *log.Logger, requestReader common.RequestReader, responseWriter common.ResponseWriter, validatorFactory common.ValidatorFactory) *API {
	return &API{
		logger:           logger,
		requestReader:    requestReader,
		responseWriter:   responseWriter,
		validatorFactory: validatorFactory,
	}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	type noteCreate struct {
		Content           string    `json:"content"`
		ExpiresAt         time.Time `json:"expiresAt"`
		ExpiresAtTimezone string    `json:"expiresAtTimezone"`
		ExpiresIn         int       `json:"expiresIn"`
		KeyHash           string    `json:"keyHash"`
	}
	note := noteCreate{}
	if err := api.requestReader.Read(r.Body, &note); err != nil {
		api.responseWriter.WriteBadRequest(w, r, err)
		return
	}

	v := api.validatorFactory()
	v.CheckField("content", note.Content != "", "must not be empty")
	v.CheckField("content", len(note.Content) <= 1024*1024, "must not be bigger than 1 MB")
	v.CheckField("keyHash", len(note.KeyHash) == 64, "must be exactly 64 bytes long")

	var t time.Time
	isExpiresInSet := note.ExpiresIn != 0
	isExpiresAtSet := note.ExpiresAt != t && note.ExpiresAtTimezone != ""
	expirationDateConflict := isExpiresInSet && isExpiresAtSet
	if expirationDateConflict || !(isExpiresInSet || isExpiresAtSet) {
		v.AddNonFieldError("either expiresIn or both expiresAt and expiresAtTimezone should be provided")
	} else {
		if note.ExpiresIn == 0 {
			v.CheckField("expiresAt", note.ExpiresAt == t, "should be provided")
			_, offset := note.ExpiresAt.Zone()
			v.CheckField("expiresAt", offset != 0, "should not include a time zone")
			v.CheckField("expiresAtTimezone", note.ExpiresAtTimezone == "", "should be provided")

			tz, err := time.LoadLocation(note.ExpiresAtTimezone)
			v.CheckField("expiresAtTimezone", err != nil, "should be a valid time zone id")

			if v.CheckIsValid() {
				expirationTime := note.ExpiresAt.In(tz)
				expiresAtGELimit := time.Until(expirationTime) >= 10*time.Minute
				v.CheckField("expiresAt", expiresAtGELimit, "should be at least 10 minutes in the future")
			}
		} else {
			expiresInGELimit := note.ExpiresIn*int(time.Second) >= int(10*time.Minute)
			v.CheckField("expiresIn", expiresInGELimit, "should not be less than 600 seconds")
		}
	}

	if !v.CheckIsValid() {
		errors := common.MergeValidationErrors(v)
		api.responseWriter.WriteValidationError(w, r, errors)
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
