package handler

import (
	"github.com/ellofae/go-concurrency-process/internal/controller"
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type UserHandler struct {
	logger      hclog.Logger
	userService domain.IUserUsecase
}

func NewUserHandler(userUsecase domain.IUserUsecase) controller.IHandler {
	return &PrivilageHandler{
		logger:      logger.GetLogger(),
		userUsecase: userUsecase,
	}
}

func (uh *UserHandler) Register(router *mux.Router) {

}
