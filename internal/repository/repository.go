package repository

import "remains_api/internal/domain"

type Storage interface {
	LoginUser(loginStruct domain.LoginStruct) (string, error)
	GetAll(userid string) ([]domain.Remains, error)
	GetFiltered(params domain.RemainRequest) ([]domain.Remains, error)
	GetOnlyGroup(group string, params domain.RemainRequest) ([]domain.Remains, error)
}
