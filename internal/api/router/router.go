package router

import (
	"log"
	"net/http"
)

func New(logger *log.Logger, noteRouter http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", noteRouter)
	return mux
}
