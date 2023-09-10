package usecase

import (
	"time"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/internal/utils"
	"github.com/hashicorp/go-hclog"
)

type PrivilegeService struct {
	logger hclog.Logger
	repo   domain.IPrivilegeRepository
}

func NewPrivilegeService(logger hclog.Logger, repo domain.IPrivilegeRepository) domain.IPrivilegeService {
	return &PrivilegeService{
		logger: logger,
		repo:   repo,
	}
}

func (ps *PrivilegeService) GetRecordByID(id int) (*entity.Privilege, error) {
	record, err := ps.repo.GetRecordByID(id)
	if err != nil {
		ps.logger.Error("Unable to get record by id", "id", id, "error", err)
		return nil, err
	}

	return record, nil
}

func (ps *PrivilegeService) GetAllRecords() ([]*entity.Privilege, error) {
	records, err := ps.repo.GetAllRecords()
	if err != nil {
		ps.logger.Error("Unable to get records of entities from privilege table", "error", err)
		return nil, err
	}

	return records, nil
}

func (ps *PrivilegeService) CreatePrivilege(req *dto.PrivilegeCreateDTO) error {
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

	if err := ps.repo.CreatePrivilege(entity); err != nil {
		ps.logger.Error("Unable to create a record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeService) UpdatePrivilege(id int, req *dto.PrivilegeUpdateDTO) error {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return err
	}

	if err := ps.repo.UpdatePrivilege(id, req); err != nil {
		ps.logger.Error("Unable to update the record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeService) DeletePrivilege(id int) error {
	if err := ps.repo.DeletePrivilege(id); err != nil {
		ps.logger.Error("Unable to delete the record in postgres database", "error", err)
		return err
	}

	return nil
}
