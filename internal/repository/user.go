package repository

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type UserRepository struct {
	logger  hclog.Logger
	storage *Storage
}

func NewUserRepository(storage *Storage) domain.IUserRepository {
	return &UserRepository{
		logger:  logger.GetLogger(),
		storage: storage,
	}
}
