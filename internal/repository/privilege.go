package repository

import (
	"context"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeRepository struct {
	logger  hclog.Logger
	storage *Storage
}

func NewPrivilegeRepository(storage *Storage) domain.IPrivilegeRepository {
	return &PrivilegeRepository{
		logger:  logger.GetLogger(),
		storage: storage,
	}
}

func (pr *PrivilegeRepository) GetRecordByID(ctx context.Context, id int) (*entity.Privilege, error) {
	query := `SELECT * FROM privileges WHERE id = $1`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	entity := &entity.Privilege{}
	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if err := tx.QueryRow(ctx, query, id).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		return nil, err
	}

	return entity, nil
}

func (pr *PrivilegeRepository) GetAllRecords(ctx context.Context) ([]*entity.Privilege, error) {
	query := `SELECT * FROM privileges`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	entities := []*entity.Privilege{}

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		privilege := &entity.Privilege{}

		err := rows.Scan(&privilege.ID, &privilege.PrivilegeTitle, &privilege.CreatedAt)
		if err != nil {
			return nil, err
		}
		entities = append(entities, privilege)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (pr *PrivilegeRepository) CreatePrivilege(ctx context.Context, req *entity.Privilege) error {
	query := `INSERT INTO privileges(privilege_title, created_at) VALUES ($1, $2)`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, req.PrivilegeTitle, req.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) UpdatePrivilege(ctx context.Context, id int, req *dto.PrivilegeUpdateDTO) error {
	query := `UPDATE privileges SET privilege_title = $2 WHERE id = $1`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, id, req.PrivilegeTitle)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) DeletePrivilege(ctx context.Context, id int) error {
	query := `DELETE FROM privileges WHERE id = $1`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
