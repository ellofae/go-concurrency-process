package usecase

import (
	"context"
	"time"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/internal/utils"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeService struct {
	logger hclog.Logger
	repo   domain.IPrivilegeRepository
}

func NewPrivilegeService(repo domain.IPrivilegeRepository) domain.IPrivilegeService {
	return &PrivilegeService{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (ps *PrivilegeService) GetRecordByID(ctx context.Context, id int) (*entity.Privilege, error) {
	record, err := ps.repo.GetRecordByID(ctx, id)
	if err != nil {
		ps.logger.Error("Unable to get record by id", "id", id, "error", err)
		return nil, err
	}

	return record, nil
}

func (ps *PrivilegeService) GetAllRecords(ctx context.Context) ([]*entity.Privilege, error) {
	records, err := ps.repo.GetAllRecords(ctx)
	if err != nil {
		ps.logger.Error("Unable to get records of entities from privilege table", "error", err)
		return nil, err
	}

	return records, nil
}

func (ps *PrivilegeService) CreatePrivilege(ctx context.Context, req *dto.PrivilegeCreateDTO) error {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return err
	}

	entity := &entity.Privilege{
		PrivilegeTitle: req.PrivilegeTitle,
		CreatedAt:      time.Now(),
	}

	if err := ps.repo.CreatePrivilege(ctx, entity); err != nil {
		ps.logger.Error("Unable to create a record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeService) UpdatePrivilege(ctx context.Context, id int, req *dto.PrivilegeUpdateDTO) error {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return err
	}

	if err := ps.repo.UpdatePrivilege(ctx, id, req); err != nil {
		ps.logger.Error("Unable to update the record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeService) DeletePrivilege(ctx context.Context, id int) error {
	if err := ps.repo.DeletePrivilege(ctx, id); err != nil {
		ps.logger.Error("Unable to delete the record in postgres database", "error", err)
		return err
	}

	return nil
}
