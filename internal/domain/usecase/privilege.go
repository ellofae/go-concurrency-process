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

type PrivilegeUsecase struct {
	logger hclog.Logger
	repo   domain.IPrivilegeRepository
}

func NewPrivilegeUsecase(repo domain.IPrivilegeRepository) domain.IPrivilegeUsecase {
	return &PrivilegeUsecase{
		logger: logger.GetLogger(),
		repo:   repo,
	}
}

func (ps *PrivilegeUsecase) GetRecordByTitle(ctx context.Context, req *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error) {
	record, err := ps.repo.GetRecordByTitle(ctx, req.PrivilegeTitle)
	if err != nil {
		ps.logger.Error("Unable to get record by title", "title requested", req.PrivilegeTitle)
		return nil, err
	}

	resp := &dto.PrivilegeResponseDTO{
		ID:             record.ID,
		PrivilegeTitle: record.PrivilegeTitle,
	}

	return resp, nil
}

func (ps *PrivilegeUsecase) GetAllRecords(ctx context.Context) ([]*dto.PrivilegeResponseDTO, error) {
	records := []*dto.PrivilegeResponseDTO{}

	entities, err := ps.repo.GetAllRecords(ctx)
	if err != nil {
		ps.logger.Error("Unable to get records of entities from privilege table", "error", err)
		return nil, err
	}

	for _, entity := range entities {
		record := &dto.PrivilegeResponseDTO{
			ID:             entity.ID,
			PrivilegeTitle: entity.PrivilegeTitle,
		}

		records = append(records, record)
	}

	return records, nil
}

func (ps *PrivilegeUsecase) CreatePrivilege(ctx context.Context, req *dto.PrivilegeDTO) error {
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

func (ps *PrivilegeUsecase) UpdatePrivilege(ctx context.Context, id int, req *dto.PrivilegeUpdateDTO) error {
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

func (ps *PrivilegeUsecase) DeletePrivilege(ctx context.Context, id int) error {
	if err := ps.repo.DeletePrivilege(ctx, id); err != nil {
		ps.logger.Error("Unable to delete the record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeUsecase) AddPrivilegeToUser(context.Context, *dto.PrivilegedUserCreateDTO) error {
	return nil
}
