package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CourierServiceRepositoryMock struct {
	Mock mock.Mock
}

func (repository *CourierServiceRepositoryMock) FindByUid(uid *string) (*entity.CourierService, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		courierService := arguments.Get(0).(entity.CourierService)
		return &courierService, nil
	}
}

func (repository *CourierServiceRepositoryMock) FindByParams(limit int, page int, sort string, filter map[string]interface{}) ([]entity.CourierService, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort, filter)
	return arguments.Get(0).([]entity.CourierService), arguments.Get(1).(*base.Pagination), nil
}

func (repository *CourierServiceRepositoryMock) FindAll(limit int, page int, sort string) ([]entity.CourierService, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort)
	return arguments.Get(0).([]entity.CourierService), arguments.Get(1).(*base.Pagination), nil
}

func (repository *CourierServiceRepositoryMock) CreateCourierService(courierService *entity.CourierService) (*entity.CourierService, error) {
	return courierService, nil
}

func (repository *CourierServiceRepositoryMock) UpdateCourierService(courierService *entity.CourierService) (*entity.CourierService, error) {
	return courierService, nil
}

func (repository *CourierServiceRepositoryMock) GetCourierService(courierService *entity.CourierService) (*entity.CourierService, error) {
	return courierService, nil
}

func (repository *CourierServiceRepositoryMock) Delete(uid string) error {
	return nil
}

func (repository *CourierServiceRepositoryMock) Update(uid string, input map[string]interface{}) error {
	return nil
}

func (repository *CourierServiceRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}

func (repository *CourierServiceRepositoryMock) CheckExistsByCourierIdShippingCode(courierUId string, shippingCode string) (bool, error) {
	arguments := repository.Mock.Called(courierUId, shippingCode)
	return arguments.Get(0).(bool), nil
}

func (repository *CourierServiceRepositoryMock) CheckExistsByUIdCourierIdShippingCode(uid string, courierUId string, shippingCode string) (bool, error) {
	return true, nil
}
