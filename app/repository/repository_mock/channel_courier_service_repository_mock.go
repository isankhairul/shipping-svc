package repository_mock

import (
	"go-klikdokter/app/model/entity"

	"github.com/stretchr/testify/mock"
)

type ChannelCourierServiceRepositoryMock struct {
	Mock mock.Mock
}

func (r *ChannelCourierServiceRepositoryMock) GetChannelCourierServicesByCourierServiceUIds(uids []*string) ([]*entity.ChannelCourierService, error) {
	panic("Not implemented")
}

func (r *ChannelCourierServiceRepositoryMock) CreateChannelCourierService(courier *entity.Courier, channel *entity.Channel, cs *entity.CourierService, priceInternal float64, status int) (*entity.ChannelCourierService, error) {
	arguments := r.Mock.Called("CreateChannelCourierService")
	f := arguments.Get(0)
	return f.(*entity.ChannelCourierService), nil
}
func (r *ChannelCourierServiceRepositoryMock) DeleteChannelCourierServiceByID(id uint64) error {
	panic("Not implemented")
}

func (r *ChannelCourierServiceRepositoryMock) DeleteChannelCourierServicesByChannelID(channelID uint64, courierID uint64) error {
	arguments := r.Mock.Called("DeleteChannelCourierServicesByChannelID")
	f := arguments.Get(0)
	if f != nil {
		return f.(error)
	}
	return nil
}
