package repository

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
)

type PrivilegeRepository struct {
	conn *postgres.PostgresConn
}

func NewPrivilegeRepository(conn *postgres.PostgresConn) domain.IPrivilegeRepository {
	return &PrivilegeRepository{
		conn: conn,
	}
}
