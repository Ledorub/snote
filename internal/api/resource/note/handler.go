package note

import (
	"errors"
	"github.com/ledorub/snote-api/internal"
	"github.com/ledorub/snote-api/internal/api/common"
	"github.com/ledorub/snote-api/internal/service"
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

func NewAPI(
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
	noteID := r.PathValue("noteID")
	keyHash := r.URL.Query().Get("key_hash")

	v := api.validatorFactory()
	v.Check(noteID != "", "note ID must not be empty")
	v.Check(keyHash != "", "key hash must not be empty")
	if !v.CheckIsValid() {
		var validationErrors []error
		for _, err := range v.GetErrors() {
			validationErrors = append(validationErrors, err)
		}
		api.responseWriter.WriteValidationError(w, r, validationErrors)
		return
	}

	note, err := api.noteService.GetNote(r.Context(), noteID, keyHash)
	if err != nil {
		if errors.Is(err, service.ErrDoesNotExist) {
			api.responseWriter.WriteNotFound(w, r)
		} else {
			api.responseWriter.WriteServerError(w, r, err)
		}
		return
	}
	noteResponse := &struct {
		ID                string    `json:"id"`
		ExpiresAt         time.Time `json:"expiresAt"`
		ExpiresAtTimeZone string    `json:"expiresAtTimeZone"`
		KeyHash           string    `json:"keyHash"`
	}{
		ID:                note.ID,
		ExpiresAt:         note.ExpiresAt,
		ExpiresAtTimeZone: note.ExpiresAtTimeZone.String(),
		KeyHash:           string(note.KeyHash),
	}
	api.responseWriter.Write(w, r, http.StatusOK, noteResponse)
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	api.responseWriter.Write(w, r, http.StatusNoContent, nil)
}
