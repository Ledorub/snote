package note

import (
	"errors"
	"github.com/ledorub/snote-api/api/common"
	"github.com/ledorub/snote-api/internal"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
	"net/http"
	"time"
)

type API struct {
	logger           *log.Logger
	requestReader    common.RequestReader
	responseWriter   common.ResponseWriter
	validatorFactory common.ValidatorFactory
	noteService      common.NoteService
}

func New(
	logger *log.Logger,
	requestReader common.RequestReader,
	responseWriter common.ResponseWriter,
	validatorFactory common.ValidatorFactory,
	noteService common.NoteService,
) *API {
	return &API{
		logger:           logger,
		requestReader:    requestReader,
		responseWriter:   responseWriter,
		validatorFactory: validatorFactory,
		noteService:      noteService,
	}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	type noteCreate struct {
		Content           string        `json:"content"`
		ExpiresAt         time.Time     `json:"expiresAt"`
		ExpiresAtTimezone string        `json:"expiresAtTimezone"`
		ExpiresIn         time.Duration `json:"expiresIn"`
		KeyHash           []byte        `json:"keyHash"`
	}
	noteData := noteCreate{}
	if err := api.requestReader.Read(r.Body, &noteData); err != nil {
		api.responseWriter.WriteBadRequest(w, r, err)
		return
	}

	v := api.validatorFactory()
	v.Check(noteData.Content != "", "content must not be empty")
	v.Check(len(noteData.KeyHash) != 0, "key hash must not be empty")

	isExpiresInSet := noteData.ExpiresIn != 0
	isExpiresAtSet := !noteData.ExpiresAt.IsZero() && noteData.ExpiresAtTimezone != ""
	expirationDateConflict := isExpiresInSet && isExpiresAtSet
	v.Check(
		expirationDateConflict || !(isExpiresInSet || isExpiresAtSet),
		"either expiresIn or both expiresAt and expiresAtTimezone should be provided",
	)
	if !v.CheckIsValid() {
		var validationErrors []error
		for _, err := range v.GetErrors() {
			validationErrors = append(validationErrors, err)
		}
		api.responseWriter.WriteValidationError(w, r, validationErrors)
		return
	}

	note, err := internal.NewNote(
		&noteData.Content,
		noteData.ExpiresIn,
		noteData.ExpiresAt,
		noteData.ExpiresAtTimezone,
		noteData.KeyHash,
	)
	if err != nil {
		api.responseWriter.WriteValidationError(w, r, []error{err})
		return
	}

	note, err = api.noteService.CreateNote(r.Context(), note)
	if err != nil {
		var validationError validator.ValidationError
		if errors.As(err, &validationError) {
			api.responseWriter.WriteValidationError(w, r, []error{validationError})
			return
		}
		api.responseWriter.WriteServerError(w, r, err)
		return
	}

	type noteCreateResponse struct {
		ID                string    `json:"id"`
		ExpiresAt         time.Time `json:"expiresAt"`
		ExpiresAtTimeZone string    `json:"expiresAtTimeZone"`
		KeyHash           []byte    `json:"keyHash"`
	}
	noteResponse := noteCreateResponse{
		ID:                note.ID,
		ExpiresAt:         note.ExpiresAt,
		ExpiresAtTimeZone: note.ExpiresAtTimeZone.String(),
		KeyHash:           note.KeyHash,
	}

	api.responseWriter.Write(w, r, http.StatusCreated, noteResponse)
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
