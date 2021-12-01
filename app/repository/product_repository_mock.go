package repository

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	Mock mock.Mock
}

func (repository *ProductRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Product, *base.Pagination, error) {
	return nil, nil, nil
}

func (repository *ProductRepositoryMock) FindByUid(uid *string) (*entity.Product, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		product := arguments.Get(0).(entity.Product)
		return &product, nil
	}
}

func (repository *ProductRepositoryMock) Create(product *entity.Product) (*entity.Product, error) {
	return product, nil
}

func (repository *ProductRepositoryMock) Delete(uid string) error {
	return nil
}

func (repository *ProductRepositoryMock) Update(uid string, input map[string]interface{}) error {
	return nil
}
