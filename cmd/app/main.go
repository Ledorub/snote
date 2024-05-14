package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ledorub/snote-api/internal/api/common"
	"github.com/ledorub/snote-api/internal/api/request"
	"github.com/ledorub/snote-api/internal/api/resource/note"
	"github.com/ledorub/snote-api/internal/api/response"
	"github.com/ledorub/snote-api/internal/api/router"
	"github.com/ledorub/snote-api/internal/config"
	"github.com/ledorub/snote-api/internal/db"
	"github.com/ledorub/snote-api/internal/encdec"
	"github.com/ledorub/snote-api/internal/logger"
	"github.com/ledorub/snote-api/internal/service"
	"github.com/ledorub/snote-api/internal/validator"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
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

	ctx, stop := createMainContext()
	defer stop()

	dbConn, err := createDBConnection(ctx, &cfg.DB)
	if err != nil {
		lg.Fatal(err)
	}
	noteRepo := createNoteRepo(lg, dbConn)
	noteService := createNoteService(lg, noteRepo)

	api := createAPI(lg, noteService)
	if err = startServer(lg, ctx, &cfg.Server, api); err != nil {
		lg.Printf("server: %w", err)
	}
	closeDBConnection(lg, dbConn)
}

func loadConfig() (*config.Config, error) {
	cfgLoader := config.NewLoader(config.LoadArgs())
	return cfgLoader.Load()
}

func createLogger() *log.Logger {
	return logger.New()
}

func createMainContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func createDBConnection(ctx context.Context, dbConfig *config.DBConfig) (*pgxpool.Pool, error) {
	dsn := db.BuildDSN(
		dbConfig.Host.Value,
		dbConfig.Port.Value,
		dbConfig.User.Value,
		dbConfig.Password.Value.GetValue(),
		dbConfig.Name.Value,
	)
	return db.CreatePool(ctx, dsn)
}

func closeDBConnection(logger *log.Logger, conn *pgxpool.Pool) {
	logger.Println("DB: closing connection...")
	conn.Close()
	logger.Println("DB: connection closed")
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

func createServer(
	logger *log.Logger,
	serverConfig *config.ServerConfig,
	api http.Handler,
) *http.Server {
	maxBytes := 1_048_576
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.Port.Value),
		Handler: http.MaxBytesHandler(api, int64(maxBytes)),
	}
}

func startServer(
	logger *log.Logger,
	ctx context.Context,
	serverConfig *config.ServerConfig,
	api http.Handler,
) error {
	srv := createServer(logger, serverConfig, api)
	srvError := make(chan error, 1)
	go func() {
		logger.Printf("server: starting on %s", srv.Addr)
		srvError <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logger.Printf("server: shutting down...")
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-srvError:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	logger.Println("server: shut down")
	return nil
}
