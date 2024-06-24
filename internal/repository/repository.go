package repository

import "remains_api/internal/domain"

type Storage interface {
	GetAll() ([]domain.Remains, error)
	GetFiltered(params domain.RemainRequest) ([]domain.Remains, error)
	GetOnlyGroup(group string, params domain.RemainRequest) ([]domain.Remains, error)
}
