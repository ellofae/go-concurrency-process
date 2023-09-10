package domain

import (
	"github.com/ellofae/go-concurrency-process/internal/domain/entity"
	"github.com/ellofae/go-concurrency-process/internal/dto"
)

type (
	IPrivilegeService interface {
		CreatePrivilege(*dto.PrivilegeCreateDTO) error
		UpdatePrivilege(int, *dto.PrivilegeUpdateDTO) error
		DeletePrivilege(int) error
		GetAllRecords() ([]*entity.Privilege, error)
		GetRecordByID(int) (*entity.Privilege, error)
	}

	IPrivilegeRepository interface {
		CreatePrivilege(*entity.Privilege) error
		UpdatePrivilege(int, *dto.PrivilegeUpdateDTO) error
		DeletePrivilege(int) error
		GetAllRecords() ([]*entity.Privilege, error)
		GetRecordByID(int) (*entity.Privilege, error)
	}
)
