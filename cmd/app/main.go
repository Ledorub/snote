package main

import (
	"fmt"
	"github.com/ledorub/snote-api/api/router"
	"github.com/ledorub/snote-api/internal/config"
	"github.com/ledorub/snote-api/internal/logger"
	"net/http"
)

func main() {
	cfg := config.New()
	fmt.Printf("Got port %d.\n", cfg.Port)
	log := logger.New()
	log.Println("Set up logger.")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router.New(log),
	}

	log.Printf("Starting the server at %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
