package repository

import (
	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/pkg/postgres"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeRepository struct {
	logger        hclog.Logger
	postgres_conn *postgres.PostgresConn
}

func NewPrivilegeRepository(logger hclog.Logger, conn *postgres.PostgresConn) domain.IPrivilegeRepository {
	return &PrivilegeRepository{
		logger:        logger,
		postgres_conn: conn,
	}
}

func (pr *PrivilegeRepository) GetRecordByID(id int) (*entity.Privilege, error) {
	query := `SELECT * FROM privileges WHERE id = $1`

	entity := &entity.Privilege{}
	if err := pr.postgres_conn.GetDB().Get(entity, query, id); err != nil {
		return nil, err
	}

	return entity, nil
}

func (pr *PrivilegeRepository) GetAllRecords() ([]*entity.Privilege, error) {
	entities := []*entity.Privilege{}

	query := `SELECT * FROM privileges`
	if err := pr.postgres_conn.GetDB().Select(&entities, query); err != nil {
		return nil, err
	}

	return entities, nil
}

func (pr *PrivilegeRepository) CreatePrivilege(req *entity.Privilege) error {
	query := `INSERT INTO privileges(privilege_title, created_at) VALUES ($1, $2)`

	tx := pr.postgres_conn.GetDB().MustBegin()
	_, err := tx.Exec(query, req.PrivilegeTitle, req.CreatedAt)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (pr *PrivilegeRepository) UpdatePrivilege(id int, req *dto.PrivilegeUpdateDTO) error {
	query := `UPDATE privileges SET privilege_title = $2 WHERE id = $1`

	tx := pr.postgres_conn.GetDB().MustBegin()
	_, err := tx.Exec(query, id, req.PrivilegeTitle)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (pr *PrivilegeRepository) DeletePrivilege(id int) error {
	query := `DELETE FROM privileges WHERE id = $1`

	tx := pr.postgres_conn.GetDB().MustBegin()
	_, err := tx.Exec(query, id)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
