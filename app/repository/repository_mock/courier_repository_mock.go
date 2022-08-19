package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CourierRepositoryMock struct {
	Mock mock.Mock
}

func (repository *CourierRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]response.CourierListResponse, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort, filter)

	return arguments.Get(0).([]response.CourierListResponse), arguments.Get(1).(*base.Pagination), nil
}

func (repository *CourierRepositoryMock) FindByCode(code string) (*entity.Courier, error) {
	arguments := repository.Mock.Called(code)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		courier := arguments.Get(0).(entity.Courier)
		return &courier, nil
	}
}

func (repository *CourierRepositoryMock) FindByUid(uid *string) (*entity.Courier, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		if len(arguments) > 1 {
			err := arguments.Get(1)
			return nil, err.(error)
		}
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
	arguments := repository.Mock.Called(uid)
	v1 := arguments.Get(0)
	if v1 != nil {
		return v1.(error)
	}
	return nil
}

func (repository *CourierRepositoryMock) Update(uid string, input map[string]interface{}) error {
	return nil
}

func (repository *CourierRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}

func (repository *CourierRepositoryMock) IsCourierHasChild(courierID uint64) *entity.CourierHasChildFlag {
	arguments := repository.Mock.Called(courierID)
	arg := arguments.Get(0)
	if arg != nil {
		return arg.(*entity.CourierHasChildFlag)
	}
	return nil
}
