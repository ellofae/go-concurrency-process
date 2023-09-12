package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/internal/errors"
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
	getRouter.HandleFunc("/priv", ph.handlePrivilegeGetByTitle)
	getRouter.HandleFunc("/priv/user", ph.handleGetAllUsers)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/priv", ph.handlePrivilegeCreate)
	postRouter.HandleFunc("/priv/user/add", ph.handleAttachPrivilegeToUser)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/priv/{id:[0-9]+}", ph.handlePrivilegeDelete)
	deleteRouter.HandleFunc("/priv/user/{id:[0-9]+}", ph.handlePrivilegeUserDelete)
}

func (ph *PrivilageHandler) handlePrivilegeGetByTitle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegeDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := ph.privilegeUsecase.GetRecordByTitle(ctx, req)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privilege record has been found", "filter title", req.PrivilegeTitle)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(record, rw); err != nil {
		ph.logger.Error("JSON sezialisation didn't complete successfuly", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (ph *PrivilageHandler) handlePrivilegeCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegeDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := ph.privilegeUsecase.CreatePrivilege(ctx, req)
	if err != nil {
		if err == errors.ErrRecordAlreadyExists {
			ph.logger.Error("Cannot create a record because record with such name already exists", "error", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ph.logger.Error("Internal error", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("Record has been created.\n"))
}

func (ph *PrivilageHandler) handlePrivilegeDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := ph.privilegeUsecase.DeletePrivilege(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No such privilege record has been found")
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf(`{"message": "Record has been deleted", "privilege_id": %d}}.`, id)))
}

func (ph *PrivilageHandler) handleGetAllUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	records, err := ph.privilegeUsecase.GetAllUsers(ctx)
	if err != nil {
		ph.logger.Error("Couldn't get records from privilege table", "error", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(records, rw); err != nil {
		ph.logger.Error("JSON sezialisation didn't complete successfuly", "error", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (ph *PrivilageHandler) handleAttachPrivilegeToUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegedUserDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err := ph.privilegeUsecase.AddPrivilegeToUser(ctx, req)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privilege record with such title exists", "error", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		} else if err == errors.ErrRecordAlreadyExists {
			ph.logger.Error("Such privilege is already assigned to the user", "error", err)
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(`{"message": "Record has been created"}`))
}

func (ph *PrivilageHandler) handlePrivilegeUserDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := ph.privilegeUsecase.DeletePrivilegeUser(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			ph.logger.Error("No privileged user record with such id has been found", "error", err)
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		ph.logger.Error("Internal error", "message", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf(`{"message": "Record has been deleted", "deleted privileged user id": %d}`, id)))
}
