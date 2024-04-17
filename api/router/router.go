package router

import (
	"github.com/ledorub/snote-api/api/common"
	"github.com/ledorub/snote-api/api/resource/note"
	"log"
	"net/http"
)

func New(logger *log.Logger, responseWriter common.ResponseWriter) *http.ServeMux {
	mux := http.NewServeMux()

	noteAPI := note.New(logger, responseWriter)

	mux.HandleFunc("POST /{noteID}", noteAPI.Create)
	mux.HandleFunc("GET /{noteID}", noteAPI.Read)
	mux.HandleFunc("DELETE /{noteID}", noteAPI.Delete)

	return mux
}
