package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ellofae/DatabaseService/app/handler"
	"github.com/ellofae/DatabaseService/app/router"
	"github.com/ellofae/go-concurrency-process/internal/domain/usecase"
	"github.com/ellofae/go-concurrency-process/internal/repository"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
	"github.com/hashicorp/go-hclog"
)

func Run() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "go-concurrency-processing",
		Level:      hclog.LevelFromString("ERROR"),
		TimeFormat: "2006-01-02 15:04:05",
	})

	// establishing connection with postgres database
	databaseConnection, err := postgres.OpenDatabaseConnection()
	if err != nil {
		logger.Error("Unable to establish connection with postgres database", "error", err)
		os.Exit(1)
	}

	// repository
	privilegeRepository := repository.NewPrivilegeRepository(logger, databaseConnection)

	// usecase
	privilegeService := usecase.NewPrivilegeService(logger, privilegeRepository)

	// handler
	privilageHandler := handler.NewPrivilegeHandler(logger, privilegeService)

	// router initialization
	router := router.InitRouter()

	// HTTP server
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	idleTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_IDLE_TIMEOUT"))

	srv := http.Server{
		Addr:         os.Getenv("SERVER_BIND_ADDRESS"),
		Handler:      router,
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}

	go func() {
		logger.Info("Starting the server")

		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Unable to start the server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	logger.Info("Signal has been recieved:", sig)

	// gracefully shutdown the server, waiting max 5 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
}
