package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ellofae/go-concurrency-process/config"
	"github.com/ellofae/go-concurrency-process/internal/controller/handler"
	"github.com/ellofae/go-concurrency-process/internal/domain/usecase"
	"github.com/ellofae/go-concurrency-process/internal/repository"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
	"github.com/gorilla/mux"
)

func Run() {
	logger := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())
	ctx := context.Background()

	connPool := postgres.OpenPoolConnection(ctx, cfg)
	if err := connPool.Ping(ctx); err != nil {
		logger.Error("Unable to ping the database connection")
		os.Exit(1)
	}
	storage := repository.NewStorage(connPool)

	router := InitRouter(storage)
	srv := InitHTTPServer(router, cfg)

	StartServer(ctx, srv)
}

func InitRouter(storage *repository.Storage) *mux.Router {
	// repositories
	privilegeRepository := repository.NewPrivilegeRepository(storage)
	userRepository := repository.NewUserRepository(storage)

	// usecases
	privilegeUsecase := usecase.NewPrivilegeUsecase(privilegeRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)

	// handlers
	privilageHandler := handler.NewPrivilegeHandler(privilegeUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// router initialization
	router := mux.NewRouter()
	privilageHandler.Register(router)
	userHandler.Register(router)

	return router
}

func InitHTTPServer(router *mux.Router, cfg *config.Config) http.Server {
	readTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.ReadTimeout)
	writeTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.WriteTimeout)
	idleTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.IdleTimeout)

	bindAddr := cfg.Server.BindAddr

	srv := http.Server{
		Addr:         bindAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}

	return srv
}

func StartServer(ctx context.Context, srv http.Server) {
	logger := logger.GetLogger()

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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
