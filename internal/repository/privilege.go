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

func (pr *PrivilegeRepository) GetRecordByID(id int) (*entity.Privilege, error) {
	query := `SELECT * FROM privileges WHERE id = $1`

	entity := &entity.Privilege{}
	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if err := conn.QueryRow(context.Background(), query, id).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		return nil, err
	}

	return entity, nil
}

func (pr *PrivilegeRepository) GetAllRecords() ([]*entity.Privilege, error) {
	query := `SELECT * FROM privileges`

	entities := []*entity.Privilege{}

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), query)
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

	// if err := pr.postgres_conn.GetDB().Select(&entities, query); err != nil {
	// 	return nil, err
	// }

	return entities, nil
}

func (pr *PrivilegeRepository) CreatePrivilege(req *entity.Privilege) error {
	query := `INSERT INTO privileges(privilege_title, created_at) VALUES ($1, $2)`

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), query, req.PrivilegeTitle, req.CreatedAt)
	if err != nil {
		return err
	}

	// tx := pr.postgres_conn.GetDB().MustBegin()
	// _, err := tx.Exec(query, req.PrivilegeTitle, req.CreatedAt)
	// if err != nil {
	// 	return err
	// }
	// tx.Commit()

	return nil
}

func (pr *PrivilegeRepository) UpdatePrivilege(id int, req *dto.PrivilegeUpdateDTO) error {
	query := `UPDATE privileges SET privilege_title = $2 WHERE id = $1`

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), query, id, req.PrivilegeTitle)
	if err != nil {
		return err
	}

	// tx := pr.postgres_conn.GetDB().MustBegin()
	// _, err := tx.Exec(query, id, req.PrivilegeTitle)
	// if err != nil {
	// 	return err
	// }
	// tx.Commit()

	return nil
}

func (pr *PrivilegeRepository) DeletePrivilege(id int) error {
	query := `DELETE FROM privileges WHERE id = $1`

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	// tx := pr.postgres_conn.GetDB().MustBegin()
	// _, err := tx.Exec(query, id)
	// if err != nil {
	// 	return err
	// }
	// tx.Commit()

	return nil
}
