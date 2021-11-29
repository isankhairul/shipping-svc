package repository

import (
	"github.com/stretchr/testify/mock"
	"gokit_example/app/model/entity"
	"gokit_example/app/model/response"
)

type ProductRepositoryMock struct {
	Mock mock.Mock
}

func (repository *ProductRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Product, *response.PaginationResponse, error) {
	arguments := repository.Mock.Called(limit, page, sort, filter)
	paginate := arguments.Get(1).(response.PaginationResponse)
	product := arguments.Get(0).([]entity.Product)

	return product, &paginate, nil
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

func (repository *ProductRepositoryMock) Update(uid *string, input map[string]interface{}) error {
	return nil
}
