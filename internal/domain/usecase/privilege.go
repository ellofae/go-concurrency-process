package usecase

import (
	"context"
	"time"

	"github.com/ellofae/go-concurrency-process/internal/domain"
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
	"github.com/ellofae/go-concurrency-process/internal/errors"
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

	_, err := ps.repo.GetRecordByTitle(ctx, req.PrivilegeTitle)
	if err == nil {
		return errors.ErrRecordAlreadyExists
	} else {
		if err != errors.ErrNoRecordFound {
			return err
		}
	}

	if err := ps.repo.CreatePrivilege(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (ps *PrivilegeUsecase) DeletePrivilege(ctx context.Context, id int) error {
	_, err := ps.GetRecordByID(ctx, id)
	if err != nil {
		if err == errors.ErrNoRecordFound {
			return errors.ErrNoRecordFound
		}
		return err
	}

	if err := ps.repo.DeletePrivilege(ctx, id); err != nil {
		return err
	}

	return nil
}

func (ps *PrivilegeUsecase) GetAllUsers(ctx context.Context) ([]*dto.PrivilegedUserDTO, error) {
	records := []*dto.PrivilegedUserDTO{}

	entities, err := ps.repo.GetAllUsers(ctx)
	if err != nil {
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

func (ps *PrivilegeUsecase) AddPrivilegeToUser(ctx context.Context, req *dto.PrivilegedUserCreateDTO) (string, error) {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return "", err
	}

	userID := req.UserID
	for _, privilege := range req.PrivilegeList {
		entity, err := ps.repo.GetRecordByTitle(ctx, privilege)
		if err != nil {
			return "", err
		}

		privileged_ids, err := ps.repo.GetUserPrivilegesByID(ctx, req.UserID)
		if err != nil {
			return "", err
		}

		for _, privilege_id := range privileged_ids {
			record_privilege, err := ps.repo.GetRecordByID(ctx, privilege_id)
			if err != nil {
				return "", err
			}

			if record_privilege == privilege {
				return privilege, errors.ErrRecordAlreadyExists
			}
		}

		if err := ps.repo.AddPrivilegeToUser(ctx, userID, entity.ID); err != nil {
			return "", err
		}

	}

	return "", nil
}

func (ps *PrivilegeUsecase) RemoveUserPrivilege(ctx context.Context, req *dto.PrivilegedUserDeleteDTO) (string, error) {
	validate := utils.NewValidator()

	if err := validate.Struct(req); err != nil {
		validation_errors := utils.ValidatorErrors(err)
		for _, error := range validation_errors {
			ps.logger.Error("Validation error", "error", error)
		}

		return "", err
	}

	userID := req.UserID

	flag := false
	for _, privilege := range req.PrivilegeList {
		entity, err := ps.repo.GetRecordByTitle(ctx, privilege)
		if err != nil {
			return "", err
		}

		privileged_ids, err := ps.repo.GetUserPrivilegesByID(ctx, req.UserID)
		if err != nil {
			return "", err
		}

		for _, privilege_id := range privileged_ids {
			record_privilege, err := ps.repo.GetRecordByID(ctx, privilege_id)
			if err != nil {
				return "", err
			}

			if record_privilege == privilege {
				flag = true
				break
			}
		}

		if flag {
			if err := ps.repo.RemoveUserPrivilege(ctx, userID, entity.ID); err != nil {
				return "", err
			}
		} else {
			return privilege, errors.ErrNoRecordFound
		}

		flag = false
	}

	return "", nil
}

func (ps *PrivilegeUsecase) DeletePrivilegeUser(ctx context.Context, id int) error {
	_, err := ps.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := ps.repo.DeletePrivilegeUser(ctx, id); err != nil {
		return err
	}

	return nil
}
