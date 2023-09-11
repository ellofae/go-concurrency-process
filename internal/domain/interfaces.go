package domain

import (
	"context"

	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
)

type (
	IPrivilegeUsecase interface {
		GetAllRecords(context.Context) ([]*dto.PrivilegeResponseDTO, error)
		GetRecordByTitle(context.Context, *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error)
		CreatePrivilege(context.Context, *dto.PrivilegeDTO) error
		DeletePrivilege(context.Context, int) error

		UpdatePrivilege(context.Context, int, *dto.PrivilegeUpdateDTO) error
	}

	IPrivilegeRepository interface {
		GetAllRecords(context.Context) ([]*entity.Privilege, error)
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
		CreatePrivilege(context.Context, *entity.Privilege) error
		DeletePrivilege(context.Context, int) error

		UpdatePrivilege(context.Context, int, *dto.PrivilegeUpdateDTO) error
	}

	IUserUsecase interface {
	}

	IUserRepository interface {
	}
)
