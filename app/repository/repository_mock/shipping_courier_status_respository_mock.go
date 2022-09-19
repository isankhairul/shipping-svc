package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type ShippingCourierStatusRepositoryMock struct {
	Mock mock.Mock
}

func (r *ShippingCourierStatusRepositoryMock) FindByParams(limit int, page int, sort string, filters map[string]interface{}) ([]entity.ShippingCourierStatus, *base.Pagination, error) {
	arguments := r.Mock.Called(limit, page, sort, filters)
	return arguments.Get(0).([]entity.ShippingCourierStatus), arguments.Get(1).(*base.Pagination), nil
}

func (r *ShippingCourierStatusRepositoryMock) FindByCode(channelID, courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.ShippingCourierStatus), nil
}

func (r *ShippingCourierStatusRepositoryMock) FindByCourierStatus(courierID uint64, statusCode string) (*entity.ShippingCourierStatus, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.ShippingCourierStatus), nil
}
