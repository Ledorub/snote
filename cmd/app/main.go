package main

import (
	"fmt"
	"github.com/ledorub/snote-api/api/common"
	"github.com/ledorub/snote-api/api/router"
	"github.com/ledorub/snote-api/internal/config"
	"github.com/ledorub/snote-api/internal/encdec"
	"github.com/ledorub/snote-api/internal/logger"
	"github.com/ledorub/snote-api/internal/request"
	"github.com/ledorub/snote-api/internal/response"
	"github.com/ledorub/snote-api/internal/validator"
	"net/http"
)

func main() {
	cfg := config.New()
	fmt.Printf("Got port %d.\n", cfg.Port)
	log := logger.New()
	log.Println("Set up logger.")

	jsonRequestReader := request.NewJSONReader(log, encdec.NewJSONDecoder())
	jsonResponseWriter := response.NewJSONWriter(log, encdec.NewJSONEncoder())
	validatorFactory := func() common.Validator { return validator.New() }
	noteAPI := router.New(log, jsonRequestReader, jsonResponseWriter, validatorFactory)

	maxBytes := 1_048_576
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: http.MaxBytesHandler(http.Handler(noteAPI), int64(maxBytes)),
	}

	log.Printf("Starting the server at %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
