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
