package usecase

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type UserUsecase struct {
	logger hclog.Logger
	repo   domain.IUserRepository
}

func NewUserUsecase(repo domain.IUserRepository) domain.IUserUsecase {
	return &UserUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}
