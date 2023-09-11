package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/internal/utils"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type PrivilageHandler struct {
	logger           hclog.Logger
	privilegeUsecase domain.IPrivilegeUsecase
}

func NewPrivilegeHandler(privilegeUsecase domain.IPrivilegeUsecase) controller.IHandler {
	return &PrivilageHandler{
		logger:           logger.GetLogger(),
		privilegeUsecase: privilegeUsecase,
	}
}

func (ph *PrivilageHandler) Register(router *mux.Router) {
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/priv/{id:[0-9]+}", ph.handlePrivilageGetByID)
	getRouter.HandleFunc("/priv", ph.handlePrivilageGetAll)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/priv", ph.handlePrivilageCreate)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/priv/{id:[0-9]+}", ph.handlePrivilageUpdate)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/priv/{id:[0-9]+}", ph.handlePrivilageDelete)
}

func (ph *PrivilageHandler) handlePrivilageGetByID(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	record, err := ph.privilegeUsecase.GetRecordByID(ctx, id)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to get the privilege record with id %d", id), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(record, rw); err != nil {
		http.Error(rw, "Unable to serialize the privilege record", http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilageHandler) handlePrivilageGetAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	records, err := ph.privilegeUsecase.GetAllRecords(ctx)
	if err != nil {
		http.Error(rw, "Unable to get the privilege records", http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(records, rw); err != nil {
		http.Error(rw, "Unable to serialize the privilege entities", http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilageHandler) handlePrivilageCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	req := &dto.PrivilegeCreateDTO{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		ph.logger.Error("Unable to decode the request data", "error", err)
		http.Error(rw, "Incorrect request has been made", http.StatusBadRequest)
		return
	}

	if err := ph.privilegeUsecase.CreatePrivilege(ctx, req); err != nil {
		http.Error(rw, "Unable to create a new privilege record", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (ph *PrivilageHandler) handlePrivilageUpdate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	req := &dto.PrivilegeUpdateDTO{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		ph.logger.Error("Unable to decode the request data", "error", err)
		http.Error(rw, "Incorrect request has been made", http.StatusBadRequest)
		return
	}

	if err := ph.privilegeUsecase.UpdatePrivilege(ctx, id, req); err != nil {
		http.Error(rw, fmt.Sprintf("Unable to update the privilege record with id %d", id), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (ph *PrivilageHandler) handlePrivilageDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := ph.privilegeUsecase.DeletePrivilege(ctx, id); err != nil {
		http.Error(rw, fmt.Sprintf("Unable to delete the privilege record with id %d", id), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
