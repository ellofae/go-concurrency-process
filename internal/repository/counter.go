package repository

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type CounterRepository struct {
	logger  hclog.Logger
	storage map[string]int
}

func NewCounterRepository() domain.ICounterRepository {
	return &CounterRepository{
		logger:  logger.GetLogger(),
		storage: make(map[string]int),
	}
}

func (cr *CounterRepository) ProcessConcurrency(name string, val int) {

}

func (cr *CounterRepository) GetStorage() map[string]int {
	return cr.storage
}

func (cr *CounterRepository) SetValue(name string, val int) int {
	cr.storage[name] = val
	return cr.storage[name]
}

func (cr *CounterRepository) IncreaseCounter(name string, val int) int {
	cr.ProcessConcurrency(name, val)
	return cr.storage[name]
}

func (cr *CounterRepository) DecreaseCounter(name string, val int) int {
	cr.ProcessConcurrency(name, val)
	return cr.storage[name]
}
