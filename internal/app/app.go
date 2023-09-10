package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ellofae/go-concurrency-process/internal/controller/handler"
	"github.com/ellofae/go-concurrency-process/internal/domain/usecase"
	"github.com/ellofae/go-concurrency-process/internal/repository"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load(".env")
	logger := hclog.Default()

	// establishing connection with postgres database
	databaseConnection, err := postgres.OpenDatabaseConnection()
	if err != nil {
		logger.Error("Unable to establish connection with postgres database", "error", err)
		os.Exit(1)
	}

	// repository
	privilegeRepository := repository.NewPrivilegeRepository(databaseConnection)

	// usecase
	privilegeService := usecase.NewPrivilegeService(privilegeRepository)

	// handler
	privilageHandler := handler.NewPrivilegeHandler(logger, privilegeService)

	// router initialization
	router := mux.NewRouter()
	privilageHandler.Register(router)

	// HTTP server
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	idleTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_IDLE_TIMEOUT"))

	srv := http.Server{
		Addr:         os.Getenv("SERVER_BIND_ADDRESS"),
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}

	go func() {
		logger.Info("Starting server...")
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error("Server was stopped", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	signal := <-sigChan
	logger.Info("signal has been recieved", "signal", signal)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
