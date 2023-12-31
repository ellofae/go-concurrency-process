package repository

import (
	"context"
	"time"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/errors"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5"
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

// func (pr *PrivilegeRepository) beginTransaction(ctx context.Context, conn *pgxpool.Conn) (pgx.Tx, error) {
// 	tx, err := conn.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tx, nil
// }

// func (pr *PrivilegeRepository) closeTransaction(ctx context.Context, tx pgx.Tx, err *error) {
// 	if *err != nil {
// 		if rollBackError := tx.Rollback(ctx); rollBackError != nil {
// 			pr.logger.Warn("Error rolling back transaction", "error", rollBackError)
// 			*err = rollBackError
// 		}
// 	} else {
// 		if commitError := tx.Commit(ctx); commitError != nil {
// 			pr.logger.Warn("Error commiting transaction", "error", commitError)
// 			*err = commitError
// 		}
// 	}
// }

func (pr *PrivilegeRepository) GetRecordByID(ctx context.Context, priv_id int) (string, error) {
	query := `SELECT * FROM privileges WHERE id = $1`
	entity := &entity.Privilege{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if err := tx.QueryRow(ctx, query, priv_id).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.ErrNoRecordFound
		}
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return entity.PrivilegeTitle, nil
}

func (pr *PrivilegeRepository) GetUserByID(ctx context.Context, user_id int) (int, error) {
	query := `SELECT * FROM privileged_users WHERE user_id = $1`
	entity := &entity.PrivilegedUser{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return -1, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	if err := tx.QueryRow(ctx, query, user_id).Scan(&entity.UserID, &entity.PrivilegeID, &entity.AssignedAt); err != nil {
		if err == pgx.ErrNoRows {
			return -1, errors.ErrNoRecordFound
		}
		return entity.UserID, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	return entity.UserID, nil
}

func (pr *PrivilegeRepository) GetRecordByTitle(ctx context.Context, title string) (*entity.Privilege, error) {
	query := `SELECT * FROM privileges WHERE privilege_title = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, query, title).Scan(&entity.ID, &entity.PrivilegeTitle, &entity.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNoRecordFound
		}
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pr *PrivilegeRepository) GetUserPrivilegesByID(ctx context.Context, user_id int) ([]int, error) {
	query := `SELECT * FROM privileged_users WHERE user_id = $1`
	privileges := []int{}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pr.storage.GetPgConnPool().Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		privilege := &entity.PrivilegedUser{}

		err := rows.Scan(&privilege.UserID, &privilege.PrivilegeID, &privilege.AssignedAt)
		if err != nil {
			return nil, err
		}
		privileges = append(privileges, privilege.PrivilegeID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return privileges, nil
}

func (pr *PrivilegeRepository) GetAllUsers(ctx context.Context) ([]*entity.PrivilegedUser, error) {
	query := `SELECT * FROM privileged_users`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entities := []*entity.PrivilegedUser{}

	conn, err := pr.storage.GetPgConnPool().Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		privilege := &entity.PrivilegedUser{}

		err := rows.Scan(&privilege.UserID, &privilege.PrivilegeID, &privilege.AssignedAt)
		if err != nil {
			return nil, err
		}
		entities = append(entities, privilege)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (pr *PrivilegeRepository) CreatePrivilege(ctx context.Context, req *entity.Privilege) error {
	query := `INSERT INTO privileges(privilege_title, created_at) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, req.PrivilegeTitle, req.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) DeletePrivilege(ctx context.Context, id int) error {
	query := `DELETE FROM privileges WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) AddPrivilegeToUser(ctx context.Context, user_id int, priv_id int) error {
	query := `INSERT INTO privileged_users (user_id, privilege_id, assigned_at) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, user_id, priv_id, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) RemoveUserPrivilege(ctx context.Context, user_id int, priv_id int) error {
	query := `DELETE FROM privileged_users WHERE user_id = $1 AND privilege_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, user_id, priv_id)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (pr *PrivilegeRepository) DeletePrivilegeUser(ctx context.Context, id int) error {
	query := `DELETE FROM privileged_users WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
