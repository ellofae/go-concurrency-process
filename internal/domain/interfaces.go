package domain

import (
	"context"

	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
)

type (
	IPrivilegeUsecase interface {
		GetRecordByTitle(context.Context, *dto.PrivilegeDTO) (*dto.PrivilegeResponseDTO, error)
		GetRecordByID(context.Context, int) (string, error)
		CreatePrivilege(context.Context, *dto.PrivilegeDTO) error
		DeletePrivilege(context.Context, int) error

		GetAllUsers(context.Context) ([]*dto.PrivilegedUserDTO, error)
		AddPrivilegeToUser(context.Context, *dto.PrivilegedUserDTO) error
	}

	IPrivilegeRepository interface {
		GetRecordByTitle(context.Context, string) (*entity.Privilege, error)
		GetRecordByID(context.Context, int) (string, error)
		CreatePrivilege(context.Context, *entity.Privilege) error
		DeletePrivilege(context.Context, int) error

		GetAllUsers(context.Context) ([]*entity.PrivilegedUser, error)
		AddPrivilegeToUser(context.Context, int, int) error
	}

	IUserUsecase interface {
	}

	IUserRepository interface {
	}
)
