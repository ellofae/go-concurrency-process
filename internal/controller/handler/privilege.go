package controller

import (
	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/hashicorp/go-hclog"
)

type PrivilageHandler struct {
	logger  hclog.Logger
	service domain.IPrivilegeService
}

func NewPrivilegeHandler(log hclog.Logger, service domain.IPrivilegeService) controller.IHandler {
	return &PrivilageHandler{
		logger:  log,
		service: service,
	}
}
