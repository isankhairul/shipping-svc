package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CourierRepositoryMock struct {
	Mock mock.Mock
}

func (repository *CourierRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.Courier, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort, filter)

	return arguments.Get(0).([]entity.Courier), arguments.Get(1).(*base.Pagination), nil
}

func (repository *CourierRepositoryMock) FindByUid(uid *string) (*entity.Courier, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		courier := arguments.Get(0).(entity.Courier)
		return &courier, nil
	}
}

func (repository *CourierRepositoryMock) CreateCourier(courier *entity.Courier) (*entity.Courier, error) {
	return courier, nil
}

func (repository *CourierRepositoryMock) Delete(uid string) error {
	return nil
}

func (repository *CourierRepositoryMock) Update(uid string, input map[string]interface{}) error {
	return nil
}

func (repository *CourierRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}
