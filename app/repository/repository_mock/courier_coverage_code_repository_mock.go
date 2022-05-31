package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type CourierCoverageCodeRepositoryMock struct {
	Mock mock.Mock
}

func (repository *CourierCoverageCodeRepositoryMock) GetCourierUid(courier *entity.Courier, uid string) error {
	return nil
}

func (repository *CourierCoverageCodeRepositoryMock) GetCourierId(courier *entity.Courier, id uint64) error {
	arguments := repository.Mock.Called(id)
	arg := arguments.Get(0).(entity.Courier)
	courier.UID = arg.UID
	return nil
}

func (repository *CourierCoverageCodeRepositoryMock) FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.CourierCoverageCode, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort)
	return arguments.Get(0).([]*entity.CourierCoverageCode), arguments.Get(1).(*base.Pagination), nil
}

func (repository *CourierCoverageCodeRepositoryMock) CombinationUnique(courierCoverageCode *entity.CourierCoverageCode, courierUid uint64, countryCode, postalCode string, id uint64) (int64, error) {
	return 0, nil
}

func (repository *CourierCoverageCodeRepositoryMock) FindByUid(uid string) (*entity.CourierCoverageCode, error) {
	arguments := repository.Mock.Called(uid)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		courierCoverageCode := arguments.Get(0).(entity.CourierCoverageCode)
		return &courierCoverageCode, nil
	}
}

func (repository *CourierCoverageCodeRepositoryMock) Create(courierCoverageCode *entity.CourierCoverageCode) (*entity.CourierCoverageCode, error) {
	return courierCoverageCode, nil
}

func (repository *CourierCoverageCodeRepositoryMock) Update(uid string, value map[string]interface{}) (*entity.CourierCoverageCode, error) {
	arguments := repository.Mock.Called(uid)
	return arguments.Get(0).(*entity.CourierCoverageCode), arguments.Error(1)
}

func (repository *CourierCoverageCodeRepositoryMock) DeleteByUid(uid string) error {
	return nil
}
