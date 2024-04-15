package note

import (
	"log"
	"net/http"
)

type API struct {
	logger *log.Logger
}

func New(logger *log.Logger) *API {
	return &API{logger: logger}
}

func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created a note."))
}

func (api *API) Read(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("A note."))
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deleted a note."))
}
