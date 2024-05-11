package note

import (
	"github.com/ledorub/snote-api/api/common"
	"log"
	"net/http"
)

func NewRouter(
	logger *log.Logger,
	requestReader common.RequestReader,
	responseWriter common.ResponseWriter,
	validatorFactory common.ValidatorFactory,
	noteService common.NoteService,
) *http.ServeMux {
	noteAPI := NewAPI(logger, requestReader, responseWriter, validatorFactory, noteService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /{noteID}", noteAPI.Create)
	mux.HandleFunc("GET /{noteID}", noteAPI.Read)
	mux.HandleFunc("DELETE /{noteID}", noteAPI.Delete)
	return mux
}
