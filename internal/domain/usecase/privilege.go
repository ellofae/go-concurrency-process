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

func (ps *PrivilegeUsecase) GetRecordByID(ctx context.Context, priv_id int) (string, error) {
	privilege, err := ps.repo.GetRecordByID(ctx, priv_id)
	if err != nil {
		ps.logger.Error("Unable to get the record title")
		return "", err
	}

	return privilege, nil
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

func (ps *PrivilegeUsecase) DeletePrivilege(ctx context.Context, id int) error {
	if err := ps.repo.DeletePrivilege(ctx, id); err != nil {
		ps.logger.Error("Unable to delete the record in postgres database", "error", err)
		return err
	}

	return nil
}

func (ps *PrivilegeUsecase) GetAllUsers(ctx context.Context) ([]*dto.PrivilegedUserDTO, error) {
	records := []*dto.PrivilegedUserDTO{}

	entities, err := ps.repo.GetAllUsers(ctx)
	if err != nil {
		ps.logger.Error("Unable to get records of privileged users table", "error", err)
		return nil, err
	}

	for _, entity := range entities {
		priv_title, err := ps.GetRecordByID(ctx, entity.PrivilegeID)
		if err != nil {
			return nil, err
		}

		record := &dto.PrivilegedUserDTO{
			UserID:    entity.UserID,
			Privilege: priv_title,
		}

		records = append(records, record)
	}

	return records, nil
}

func (ps *PrivilegeUsecase) AddPrivilegeToUser(ctx context.Context, req *dto.PrivilegedUserDTO) error {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return err
	}

	entity, err := ps.repo.GetRecordByTitle(ctx, req.Privilege)
	if err != nil {
		ps.logger.Error("Unable to get record by title", "title requested", req.Privilege)
		return err
	}

	if err := ps.repo.AddPrivilegeToUser(ctx, req.UserID, entity.ID); err != nil {
		ps.logger.Error("Unable to add privilege to user", "error", err)
		return err
	}

	return nil
}
