package domain

import (
	"context"

	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
)

type (
	IPrivilegeUsecase interface {
		CreatePrivilege(context.Context, *dto.PrivilegeCreateDTO) error
		UpdatePrivilege(context.Context, int, *dto.PrivilegeUpdateDTO) error
		DeletePrivilege(context.Context, int) error
		GetAllRecords(context.Context) ([]*entity.Privilege, error)
		GetRecordByID(context.Context, int) (*entity.Privilege, error)
	}

	IPrivilegeRepository interface {
		CreatePrivilege(context.Context, *entity.Privilege) error
		UpdatePrivilege(context.Context, int, *dto.PrivilegeUpdateDTO) error
		DeletePrivilege(context.Context, int) error
		GetAllRecords(context.Context) ([]*entity.Privilege, error)
		GetRecordByID(context.Context, int) (*entity.Privilege, error)
	}

	IUserUsecase interface {
	}

	IUserRepository interface {
	}
)
