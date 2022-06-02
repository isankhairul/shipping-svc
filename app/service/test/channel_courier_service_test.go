package test

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

func TestCreateChannelCourier(t *testing.T) {
	input := request.SaveChannelCourierRequest{
		CourierUID:   "courier_1",
		ChannelUID:   "channel_1",
		PrioritySort: 10,
		Status:       1,
		CourierServiceUIDs: []*request.CourierServiceDTO{
			{
				PriceInternal:     9.4,
				Status:            1,
				CourierServiceUid: "courier_service_1",
			},
		},
	}
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	//return nil when creating
	channelCourierRepo.Mock.On("GetChannelCourierByIds", mock.Anything, mock.Anything).Return(nil)
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(entity.CourierService{Status: 1})
	channelCourierServiceRepo.Mock.
		On("CreateChannelCourierService", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&entity.ChannelCourierService{Courier: courier, Channel: channel})

	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier, Channel: channel,
	}
	channelCourierRepo.Mock.On("CreateChannelCourier", mock.Anything).Return(cc)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", "123").Return(cc)
	result, msg := channelCourierService.CreateChannelCourier(input)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
	assert.Equal(t, "123", result.Uid, "Uid should be 123")
	assert.Equal(t, "Courier 1", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "Channel 1", result.ChannelName, "CourierName must be test name")
}


func TestUpdateChannelCourier(t *testing.T) {
	input := request.UpdateChannelCourierRequest{
		PrioritySort: 10,
		Status:       1,
		CourierServiceUIDs: []*request.CourierServiceDTO{
			{
				PriceInternal:     7.4,
				Status:            1,
				CourierServiceUid: "courier_service_1_1",
			},
			{
				PriceInternal:     8.4,
				Status:            1,
				CourierServiceUid: "courier_service_1_2",
			},
		},
	}
	input.Uid = "123"
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1", ID: 1}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1", ID: 1}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	//return nil when creating
	channelCourierRepo.Mock.On("GetChannelCourierByIds", mock.Anything, mock.Anything).Return(nil)
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(entity.CourierService{Status: 1})
	channelCourierServiceRepo.Mock.
		On("CreateChannelCourierService", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&entity.ChannelCourierService{Courier: courier, Channel: channel, PriceInternal: 7})

	channelCourierServiceRepo.Mock.On("DeleteChannelCourierServiceByID", mock.Anything).Return(nil)

	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier, Channel: channel,
		ChannelCourierServices: []*entity.ChannelCourierService{
			{PriceInternal: 7, CourierServiceID: 1, CourierID: 1, CourierService: &entity.CourierService{
				BaseIDModel: base.BaseIDModel{ID: 2, UID: "courier_service_1_2"},
			}},
		},
	}
	channelCourierRepo.Mock.On("GetChannelCourierByUID", "123").Return(cc)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", "123").Return(cc)
	result, msg := channelCourierService.UpdateChannelCourier(input)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
	assert.Equal(t, "123", result.Uid, "Uid should be 123")
	assert.Equal(t, "Courier 1", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "Channel 1", result.ChannelName, "CourierName must be test name")
}


func TestListChannelCouriers(t *testing.T) {
	req := request.ChannelCourierListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1", ID: 1}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1", ID: 1}, ChannelName: "Channel 1"}

	courier2 := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_2", ID: 2}, CourierName: "Courier 2"}
	channel2 := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_2", ID: 2}, ChannelName: "Channel 2"}

	items := []*entity.ChannelCourier{
		{
			Courier: courier, Channel: channel, PrioritySort: 4, Status: 1,
			ChannelCourierServices: []*entity.ChannelCourierService{
				{Courier: courier, Channel: channel, PriceInternal: 11, Status: 1,
					BaseIDModel: base.BaseIDModel{UID: "courier_service_1_1"},
					CourierService: &entity.CourierService{
						ShippingName: "courier service name 11",
					}},
				{Courier: courier, Channel: channel, PriceInternal: 12, Status: 1,
					BaseIDModel: base.BaseIDModel{UID: "courier_service_1_2"},
					CourierService: &entity.CourierService{
						ShippingName: "courier service name 12",
					}},
			},
		},
		{
			Courier: courier2, Channel: channel2, PrioritySort: 7, Status: 1,
			ChannelCourierServices: []*entity.ChannelCourierService{
				{
					Courier: courier2, Channel: channel2, PriceInternal: 22,
					Status: 1, BaseIDModel: base.BaseIDModel{UID: "courier_service_2_1"}, CourierService: &entity.CourierService{
						ShippingName: "courier service name 21",
					}},
				{
					Courier: courier2, Channel: channel2, PriceInternal: 23, Status: 1,
					BaseIDModel: base.BaseIDModel{UID: "courier_service_2_2"},
					CourierService: &entity.CourierService{
						ShippingName: "courier service name 22",
					}},
			},
		},
	}

	paginationResult := base.Pagination{
		Records:   2,
		Limit:     10,
		Page:      1,
		TotalPage: 1,
	}

	channelCourierRepo.Mock.On("FindByPagination", 10, 1, "", mock.Anything).Return(items, &paginationResult)
	results, pagination, msg := channelCourierService.ListChannelCouriers(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, len(items), len(results), "Count of cc must be 2")
	assert.Equal(t, int64(len(results)), pagination.Records, "Total record pagination must be 2")

}

func TestGetChannelCourier(t *testing.T) {
	uid := "123"
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: uid},
		Courier:     courier, Channel: channel,
	}

	channelCourierRepo.Mock.On("GetChannelCourierByUID", uid).Return(cc)
	result, msg := channelCourierService.GetChannelCourier(uid)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
	assert.Equal(t, "123", result.Uid, "Uid should be 123")
	assert.Equal(t, "Courier 1", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "Channel 1", result.ChannelName, "CourierName must be test name")

}

func TestDeleteChannelCourier(t *testing.T) {
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier,
		Channel:     channel,
	}
	cc.ID = 123

	channelCourierRepo.Mock.On("GetChannelCourierByUID", "123").Return(cc)
	channelCourierServiceRepo.Mock.On("DeleteChannelCourierServicesByChannelID", "channel_1", "courier_1").Return(nil)
	channelCourierRepo.Mock.On("DeleteChannelCourierByID", cc.ID).Return(nil)
	msg := channelCourierService.DeleteChannelCourier("123")

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}
