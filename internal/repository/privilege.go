package repository

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeRepository struct {
	logger hclog.Logger
	conn   *postgres.PostgresConn
}

func NewPrivilegeRepository(log hclog.Logger, conn *postgres.PostgresConn) domain.IPrivilegeRepository {
	return &PrivilegeRepository{
		logger: log,
		conn:   conn,
	}
}
