package usecase

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
)

type PrivilegeService struct {
	repo domain.IPrivilegeRepository
}

func NewPrivilegeService(repo domain.IPrivilegeRepository) domain.IPrivilegeService {
	return &PrivilegeService{
		repo: repo,
	}
}
