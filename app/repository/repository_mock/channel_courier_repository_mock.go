package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type ChannelCourierRepositoryMock struct {
	Mock mock.Mock
}

func (r *ChannelCourierRepositoryMock) UpdateChannelCourier(uid string, data map[string]interface{}) error {
	return nil
}

func (r *ChannelCourierRepositoryMock) CreateChannelCourier(channelCourier *entity.ChannelCourier) (*entity.ChannelCourier, error) {
	arguments := r.Mock.Called()
	f := arguments.Get(0)
	if f == nil {
		return nil, nil
	}
	return f.(*entity.ChannelCourier), nil
}

func (r *ChannelCourierRepositoryMock) GetChannelCourier(channelID int, courierID int) (*entity.ChannelCourier, error) {
	panic("Not implemented")
}

func (r *ChannelCourierRepositoryMock) GetChannelCourierByUID(uid string) (*entity.ChannelCourier, error) {
	arguments := r.Mock.Called(uid)
	f := arguments.Get(0)
	if f == nil {
		if len(arguments) > 1 {
			return nil, arguments.Get(1).(error)
		}
		return nil, nil
	}
	if len(arguments) > 1 {
		return f.(*entity.ChannelCourier), arguments.Get(1).(error)
	}
	return f.(*entity.ChannelCourier), nil
}

func (r *ChannelCourierRepositoryMock) GetChannelCourierByIds(channelID uint64, courierID uint64) (*entity.ChannelCourier, error) {
	arguments := r.Mock.Called("GetChannelCourierByIds")
	f := arguments.Get(0)
	if f == nil {
		return nil, nil
	}
	return f.(*entity.ChannelCourier), nil
}

func (r *ChannelCourierRepositoryMock) FindByPagination(limit int, page int, sort string, filters map[string]interface{}) ([]*entity.ChannelCourier, *base.Pagination, error) {
	arguments := r.Mock.Called(limit, page, sort)
	f := arguments.Get(0).([]*entity.ChannelCourier)
	f1 := arguments.Get(1).(*base.Pagination)
	return f, f1, nil
}

func (r *ChannelCourierRepositoryMock) DeleteChannelCourierByID(id uint64) error {
	return nil
}

func (r *ChannelCourierRepositoryMock) FindCourierByUID(uid string) (*entity.Courier, error) {
	arguments := r.Mock.Called("FindCourierByUID")
	c := arguments.Get(0)
	if c == nil {
		return nil, arguments.Get(1).(error)
	}
	return arguments.Get(0).(*entity.Courier), nil
}

func (r *ChannelCourierRepositoryMock) FindChannelByUID(uid string) (*entity.Channel, error) {
	arguments := r.Mock.Called("FindChannelByUID")
	c := arguments.Get(0)
	if c == nil {
		return nil, arguments.Get(1).(error)
	}
	return arguments.Get(0).(*entity.Channel), nil
}

func (r *ChannelCourierRepositoryMock) IsHasChannelCourierService(channelCourierID uint64) bool {
	arguments := r.Mock.Called()

	if arguments.Get(0) == nil {
		return false
	}
	return arguments.Get(0).(bool)
}
