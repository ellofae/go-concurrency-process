package handler

import (
	"net/http"
	"strconv"

	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Response struct {
	Message string `json:"message"`
}

type CounterHandler struct {
	logger         hclog.Logger
	counterUsecase domain.ICounterUsecase
}

func NewCounterHandler(counterUsecase domain.ICounterUsecase) controller.IHandler {
	return &CounterHandler{
		logger:         logger.GetLogger(),
		counterUsecase: counterUsecase,
	}
}

func (ch *CounterHandler) Register(router *mux.Router) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/set/{val:[0-9]+}", ch.handleSetCounter).Queries("name", "{[a-z]+}")
	getRouter.HandleFunc("/set/{val:[0-9]+}", ch.handleSetCounter)

	getRouter.HandleFunc("/inc/{val:[0-9]+}", ch.handleIncreaseCounter).Queries("name", "{[a-z]+}")
	getRouter.HandleFunc("/inc/{val:[0-9]+}", ch.handleIncreaseCounter)

	getRouter.HandleFunc("/dec/{val:[0-9]+}", ch.handleDecreaseCounter).Queries("name", "{[a-z]+}")
	getRouter.HandleFunc("/dec/{val:[0-9]+}", ch.handleDecreaseCounter)
}

func (ch *CounterHandler) handleSetCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")

	_ = ch.counterUsecase.SetValue(name, val)

	rw.WriteHeader(http.StatusOK)
}

func (ch *CounterHandler) handleIncreaseCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")

	currentValue := ch.counterUsecase.IncreaseCounter(name, val)
	if currentValue == -1 {
		ch.logger.Warn("MaxInt ceiling has been hit")
		return
	} else if currentValue == 2 {
		ch.logger.Warn("Default value is set to zero")
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (ch *CounterHandler) handleDecreaseCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")

	currentValue := ch.counterUsecase.DecreaseCounter(name, val)
	if currentValue == -1 {
		ch.logger.Warn("Current value is zero, set the new one to continue decrementing")
		return
	}

	rw.WriteHeader(http.StatusOK)
}
