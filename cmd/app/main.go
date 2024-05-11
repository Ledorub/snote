package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ledorub/snote-api/api/common"
	"github.com/ledorub/snote-api/api/resource/note"
	"github.com/ledorub/snote-api/api/router"
	"github.com/ledorub/snote-api/internal/config"
	"github.com/ledorub/snote-api/internal/db"
	"github.com/ledorub/snote-api/internal/encdec"
	"github.com/ledorub/snote-api/internal/logger"
	"github.com/ledorub/snote-api/internal/request"
	"github.com/ledorub/snote-api/internal/response"
	"github.com/ledorub/snote-api/internal/service"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
	"net/http"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	lg := createLogger()
	lg.Println("Set up logger.")

	cfgPretty, err := cfg.Pretty()
	if err != nil {
		log.Fatalf("unable to print config: %w", err)
	}
	log.Printf("Config:\n%v", cfgPretty)

	dbConn, err := createDBConnection(&cfg.DB)
	if err != nil {
		lg.Fatal(err)
	}
	noteRepo := createNoteRepo(lg, dbConn)
	noteService := createNoteService(lg, noteRepo)

	api := createAPI(lg, noteService)
	srv := createServer(lg, &cfg.Server, api)

	lg.Printf("Starting the server at %s", srv.Addr)
	err = srv.ListenAndServe()
	lg.Fatal(err)
}

func loadConfig() (*config.Config, error) {
	cfgLoader := config.NewLoader(config.LoadArgs())
	return cfgLoader.Load()
}

func createLogger() *log.Logger {
	return logger.New()
}

func createDBConnection(dbConfig *config.DBConfig) (*pgxpool.Pool, error) {
	dsn := db.BuildDSN(
		dbConfig.Host.Value,
		dbConfig.Port.Value,
		dbConfig.User.Value,
		dbConfig.Password.Value.GetValue(),
		dbConfig.Name.Value,
	)
	return db.CreatePool(context.Background(), dsn)
}

func createNoteRepo(logger *log.Logger, dbConn *pgxpool.Pool) *db.NoteRepository {
	return db.NewNoteRepository(logger, db.New(dbConn))
}

func createNoteService(logger *log.Logger, repo *db.NoteRepository) *service.NoteService {
	return service.New(logger, repo)
}

func createAPI(logger *log.Logger, service *service.NoteService) *http.ServeMux {
	jsonRequestReader := request.NewJSONReader(logger, encdec.NewJSONDecoder())
	jsonResponseWriter := response.NewJSONWriter(logger, encdec.NewJSONEncoder())
	validatorFactory := func() common.Validator { return validator.New() }
	noteAPI := note.NewRouter(logger, jsonRequestReader, jsonResponseWriter, validatorFactory, service)
	return router.New(logger, noteAPI)
}

func createServer(logger *log.Logger, serverConfig *config.ServerConfig, api http.Handler) *http.Server {
	maxBytes := 1_048_576
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port),
		Handler: http.MaxBytesHandler(api, int64(maxBytes)),
	}
}
