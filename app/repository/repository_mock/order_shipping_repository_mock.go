package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"

	"github.com/stretchr/testify/mock"
)

type OrderShippingRepositoryMock struct {
	Mock mock.Mock
}

func (r *OrderShippingRepositoryMock) Create(input *entity.OrderShipping) (*entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) Update(input *entity.OrderShipping) (*entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) Upsert(input *entity.OrderShipping) (*entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) FindByOrderNo(orderNo string) (*entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) FindByUID(uid string) (*entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) FindByParams(limit, page int, sort string, filter map[string]interface{}) ([]response.GetOrderShippingList, *base.Pagination, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 2 {
		if arguments.Get(2) != nil {
			return nil, nil, arguments.Get(2).(error)
		}
	}

	return arguments.Get(0).([]response.GetOrderShippingList), arguments.Get(1).(*base.Pagination), nil
}

func (r *OrderShippingRepositoryMock) FindByUIDs(channelUID string, uid []string) ([]entity.OrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).([]entity.OrderShipping), nil
}

func (r *OrderShippingRepositoryMock) Download(filter map[string]interface{}) ([]response.DownloadOrderShipping, error) {
	arguments := r.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).([]response.DownloadOrderShipping), nil
}
