package test

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
// var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
// var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}

//var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

func init() {
}

func TestCreateChannelCourier(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierRequest{
		CourierUID:   "courier_1",
		ChannelUID:   "channel_1",
		PrioritySort: 10,
		Status:       1,
	}
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	//return nil when creating
	channelCourierRepo.Mock.On("GetChannelCourierByIds", mock.Anything, mock.Anything).Return(nil)

	status1 := int32(1)
	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier, Channel: channel,
		Status: &status1,
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
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierRequest{
		PrioritySort: 10,
		Status:       1,
	}
	input.Uid = "123"
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1", ID: 1}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1", ID: 1}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	//return nil when creating
	channelCourierRepo.Mock.On("GetChannelCourierByIds", mock.Anything, mock.Anything).Return(nil)

	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier, Channel: channel,
		Status: &input.Status,
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
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	req := request.ChannelCourierListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1", ID: 1}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1", ID: 1}, ChannelName: "Channel 1"}

	courier2 := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_2", ID: 2}, CourierName: "Courier 2"}
	channel2 := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_2", ID: 2}, ChannelName: "Channel 2"}

	status1 := int32(1)
	items := []*entity.ChannelCourier{
		{
			Courier: courier, Channel: channel, PrioritySort: 4, Status: &status1,
		},
		{
			Courier: courier2, Channel: channel2, PrioritySort: 7, Status: &status1,
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
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	uid := "123"
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	status1 := int32(1)
	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: uid},
		Courier:     courier, Channel: channel,
		Status: &status1,
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
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	cc := &entity.ChannelCourier{
		BaseIDModel: base.BaseIDModel{UID: "123"},
		Courier:     courier,
		Channel:     channel,
	}
	cc.ID = 123

	channelCourierRepo.Mock.On("GetChannelCourierByUID", "123").Return(cc)
	channelCourierServiceRepo.Mock.On("DeleteChannelCourierServicesByChannelID", mock.Anything).Return(nil)
	channelCourierRepo.Mock.On("DeleteChannelCourierByID", cc.ID).Return(nil)
	msg := channelCourierService.DeleteChannelCourier("123")

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

//create failed test cases

func TestCreateChannelCourierFailedWithInvalidCourier(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierRequest{
		CourierUID:   "courier_1",
		ChannelUID:   "channel_1",
		PrioritySort: 10,
		Status:       1,
	}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(nil, errors.New("Not found"))
	result, msg := channelCourierService.CreateChannelCourier(input)

	assert.Nil(t, result)
	assert.Equal(t, message.ErrCourierNotFound.Message, msg.Message, "Courier not found")
	assert.Equal(t, message.ErrCourierNotFound.Code, msg.Code, "Courier not found")
}

func TestCreateChannelCourierFailedWithInvalidChannel(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierRequest{
		CourierUID:   "courier_1",
		ChannelUID:   "channel_1",
		PrioritySort: 10,
		Status:       1,
	}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(nil, errors.New("Not found"))
	result, msg := channelCourierService.CreateChannelCourier(input)

	assert.Nil(t, result)
	assert.Equal(t, message.ErrCourierNotFound.Message, msg.Message, "Courier not found")
	assert.Equal(t, message.ErrCourierNotFound.Code, msg.Code, "Courier not found")
}

func TestCreateChannelCourierFailedWithDuplicatedChannelCourier(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierRequest{
		CourierUID:   "courier_1",
		ChannelUID:   "channel_1",
		PrioritySort: 10,
		Status:       1,
	}
	courier := &entity.Courier{BaseIDModel: base.BaseIDModel{UID: "courier_1"}, CourierName: "Courier 1"}
	channel := &entity.Channel{BaseIDModel: base.BaseIDModel{UID: "channel_1"}, ChannelName: "Channel 1"}
	channelCourierRepo.Mock.On("FindCourierByUID", mock.Anything).Return(courier)
	channelCourierRepo.Mock.On("FindChannelByUID", mock.Anything).Return(channel)

	channelCourierRepo.Mock.On("GetChannelCourierByIds", mock.Anything, mock.Anything).
		Return(&entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: "dup"}, Status: &input.Status})
	result, msg := channelCourierService.CreateChannelCourier(input)

	assert.NotNil(t, result)
	assert.Equal(t, result.Uid, "dup")
	assert.Equal(t, msg.Message, message.ErrChannelCourierFound.Message, "Duplicated channel courier")
	assert.Equal(t, msg.Code, message.ErrChannelCourierFound.Code, "Duplicated channel courier")
}

/*
	Test failed cases with update
*/

func TestUpdateChannelCourierFailedWithChannelCourierNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierRequest{
		Uid:          "123",
		PrioritySort: 10,
		Status:       1,
	}
	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(nil, errors.New("Not found channel courier"))
	result, msg := channelCourierService.UpdateChannelCourier(input)

	assert.Nil(t, result)
	assert.Equal(t, message.ErrChannelCourierNotFound.Message, msg.Message, "Channel Courier not found")
	assert.Equal(t, message.ErrChannelCourierNotFound.Code, msg.Code, "Channel Courier not found")
}

func TestGetChannelCourierWithChannelCourierNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(nil)

	result, msg := channelCourierService.GetChannelCourier("123")

	assert.Nil(t, result)
	assert.Equal(t, message.ErrChannelCourierNotFound.Message, msg.Message, "ErrChannelCourierNotFound")
	assert.Equal(t, message.ErrChannelCourierNotFound.Code, msg.Code, "ErrChannelCourierNotFound")
}

func TestGetChannelCourierWithChannelCourierDbError(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(&entity.ChannelCourier{}, errors.New("Database issue"))

	result, msg := channelCourierService.GetChannelCourier("123")

	assert.Nil(t, result)
	assert.Equal(t, message.ErrDB.Message, msg.Message, "Database issue")
	assert.Equal(t, message.ErrDB.Code, msg.Code, "Database issue")
}

func TestDeleteChannelCourierFailedWithChannelCourierNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(nil, errors.New("Not found channel courier"))

	msg := channelCourierService.DeleteChannelCourier("123")

	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrChannelCourierNotFound.Message, msg.Message, "ErrUnableToDeleteChannelCourier")
	assert.Equal(t, message.ErrChannelCourierNotFound.Code, msg.Code, "ErrUnableToDeleteChannelCourier")
}

func TestDeleteChannelCourierFailedWithChannelCourierServiceIssue(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierService = service.NewChannelCourierService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(&entity.ChannelCourier{})
	channelCourierServiceRepo.Mock.On("DeleteChannelCourierServicesByChannelID", mock.Anything).Return(errors.New("Unable to delete channel courier service"))
	msg := channelCourierService.DeleteChannelCourier("123")

	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrUnableToDeleteChannelCourier.Message, msg.Message, "ErrUnableToDeleteChannelCourier")
	assert.Equal(t, message.ErrUnableToDeleteChannelCourier.Code, msg.Code, "ErrUnableToDeleteChannelCourier")
}
