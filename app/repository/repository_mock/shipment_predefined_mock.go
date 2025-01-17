package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type ShipmentPredefinedMock struct {
	Mock mock.Mock
}

func (repository *ShipmentPredefinedMock) GetShipmentPredefinedByUid(uid string) (*entity.ShippmentPredefined, error) {
	arguments := repository.Mock.Called(uid)
	ret := arguments.Get(0)
	//var s entity.ShippmentPredefined
	if ret != nil {
		s := ret.(entity.ShippmentPredefined)
		return &s, nil
	}
	if len(arguments) > 1 {
		return nil, arguments.Get(1).(error)
	}
	return nil, nil
}
func (repository *ShipmentPredefinedMock) GetAll(limit int, page int, sort string, filter map[string]interface{}) ([]*entity.ShippmentPredefined, *base.Pagination, error) {
	arguments := repository.Mock.Called(limit, page, sort)
	items := arguments.Get(0)
	return items.([]*entity.ShippmentPredefined), arguments.Get(1).(*base.Pagination), nil
}

func (repository *ShipmentPredefinedMock) UpdateShipmentPredefined(dto entity.ShippmentPredefined) (*entity.ShippmentPredefined, error) {
	arguments := repository.Mock.Called(dto)
	ret := arguments.Get(0).(entity.ShippmentPredefined)
	return &ret, nil
}

func (repository *ShipmentPredefinedMock) GetListByType(Type string) ([]entity.ShippmentPredefined, error) {
	arguments := repository.Mock.Called()

	if len(arguments) > 1 && arguments.Get(1) != nil {
		return []entity.ShippmentPredefined{}, arguments.Get(1).(error)
	}

	if arguments.Get(0) == nil {
		return []entity.ShippmentPredefined{}, nil
	}

	return arguments.Get(0).([]entity.ShippmentPredefined), nil
}
