package handler

import (
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
	getRouter.HandleFunc("/priv", ph.handlePrivilegeGetByTitle)
	getRouter.HandleFunc("/priv/user/", ph.handleGetAllUsers)

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
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(record, rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilageHandler) handlePrivilegeCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegeDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ph.privilegeUsecase.CreatePrivilege(ctx, req); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (ph *PrivilageHandler) handlePrivilegeDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := ph.privilegeUsecase.DeletePrivilege(ctx, id); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (ph *PrivilageHandler) handleGetAllUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	records, err := ph.privilegeUsecase.GetAllUsers(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = utils.ToJSON(records, rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *PrivilageHandler) handleAttachPrivilegeToUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	req := &dto.PrivilegedUserDTO{}
	if err := utils.StructDecode(r, req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ph.privilegeUsecase.AddPrivilegeToUser(ctx, req); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (ph *PrivilageHandler) handlePrivilegeUserDelete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := ph.privilegeUsecase.DeletePrivilegeUser(ctx, id); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
