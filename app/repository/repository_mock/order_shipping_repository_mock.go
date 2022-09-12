package repository_mock

import (
	"go-klikdokter/app/model/entity"

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
