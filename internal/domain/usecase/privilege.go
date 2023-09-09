package usecase

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeService struct {
	logger hclog.Logger
	repo   domain.IPrivilegeRepository
}

func NewPrivilegeService(log hclog.Logger, repo domain.IPrivilegeRepository) domain.IPrivilegeService {
	return &PrivilegeService{
		logger: log,
		repo:   repo,
	}
}
