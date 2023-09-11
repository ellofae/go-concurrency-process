package handler

import (
	"net/http"
	"strconv"

	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/utils"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

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

	currentValue := ch.counterUsecase.SetValue(name, val)
	if err := utils.ToJSON(currentValue, rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ch *CounterHandler) handleIncreaseCounter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	val, _ := strconv.Atoi(vars["val"])
	name := r.URL.Query().Get("name")

	currentValue := ch.counterUsecase.IncreaseCounter(name, val)
	if currentValue == -1 {
		ch.logger.Warn("MaxInt ceiling has been hit")
		http.Error(rw, "Ceiling of MaxInt has been hit", http.StatusBadRequest)
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
		ch.logger.Warn("0 floor has been hit")
		http.Error(rw, "Cannot reduce the value when it is is 0", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
