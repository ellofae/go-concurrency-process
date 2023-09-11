package usecase

import (
	"math"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type CounterUsecase struct {
	logger hclog.Logger
	repo   domain.ICounterRepository
}

func NewCounterUsecase(repo domain.ICounterRepository) domain.ICounterUsecase {
	return &CounterUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (cu *CounterUsecase) SetValue(name string, val int) int {
	return cu.repo.SetValue(name, val)
}

func (cu *CounterUsecase) IncreaseCounter(name string, val int) int {
	storage := cu.repo.GetStorage()
	if _, ok := storage[name]; ok {
		if storage[name]+val > math.MaxInt {
			return -1
		}
	} else {
		storage[name] = 0
	}

	return cu.repo.IncreaseCounter(name, val)
}

func (cu *CounterUsecase) DecreaseCounter(name string, val int) int {
	storage := cu.repo.GetStorage()
	if _, ok := storage[name]; ok {
		if storage[name]-val >= 0 {
			return cu.repo.DecreaseCounter(name, val)
		}
	}

	return -1
}
