package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type PrivilageHandler struct {
	logger  hclog.Logger
	service domain.IPrivilegeService
}

func NewPrivilegeHandler(logger hclog.Logger, service domain.IPrivilegeService) controller.IHandler {
	return &PrivilageHandler{
		logger:  logger,
		service: service,
	}
}

func (ph *PrivilageHandler) Register(router *mux.Router) {
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/priv", ph.handlePrivilageCreate)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/priv", ph.test)
}

func (ph *PrivilageHandler) handlePrivilageCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	req := &dto.PrivilegeCreateDTO{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		ph.logger.Error("Unable to decode the request", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ph.logger.Info("testing", "object", req)
	rw.WriteHeader(http.StatusCreated)
}

func (ph *PrivilageHandler) test(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ph.logger.Info("get test")
	rw.WriteHeader(http.StatusCreated)
}
