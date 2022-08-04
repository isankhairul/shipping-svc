package repository_mock

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/response"

	"github.com/stretchr/testify/mock"
)

type ChannelCourierServiceRepositoryMock struct {
	Mock mock.Mock
}

func (r *ChannelCourierServiceRepositoryMock) GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error) {
	panic("Not implemented")
}

func (r *ChannelCourierServiceRepositoryMock) CreateChannelCourierService(data *entity.ChannelCourierService) (*entity.ChannelCourierService, error) {
	return data, nil
}

func (r *ChannelCourierServiceRepositoryMock) DeleteChannelCourierServiceByID(id uint64) error {
	return nil
}

func (r *ChannelCourierServiceRepositoryMock) DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error {
	arguments := r.Mock.Called("DeleteChannelCourierServicesByChannelID")
	f := arguments.Get(0)
	if f != nil {
		return f.(error)
	}
	return nil
}

func (r *ChannelCourierServiceRepositoryMock) GetChannelCourierService(channelCourierID, courierServiceID uint64) (*entity.ChannelCourierService, error) {
	arguments := r.Mock.Called(channelCourierID, courierServiceID)

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.ChannelCourierService), nil
}

func (r *ChannelCourierServiceRepositoryMock) GetChannelCourierServiceByUID(uid string) (*entity.ChannelCourierService, error) {
	arguments := r.Mock.Called(uid)

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*entity.ChannelCourierService), nil
}
func (r *ChannelCourierServiceRepositoryMock) UpdateChannelCourierService(uid string, updates map[string]interface{}) error {
	return nil
}
func (r *ChannelCourierServiceRepositoryMock) FindByParams(limit, page int, sort string, filters map[string]interface{}) ([]entity.ChannelCourierService, *base.Pagination, error) {
	arguments := r.Mock.Called(limit, page, sort, filters)
	return arguments.Get(0).([]entity.ChannelCourierService), arguments.Get(1).(*base.Pagination), nil
}

func (r *ChannelCourierServiceRepositoryMock) GetChannelCourierListByChannelUID(channel_uid string, limit int, page int, sort, dir string, filter map[string][]string) ([]response.CourierServiceByChannelResponse, *base.Pagination, error) {
	arguments := r.Mock.Called()
	return arguments.Get(0).([]response.CourierServiceByChannelResponse), arguments.Get(1).(*base.Pagination), nil
}
