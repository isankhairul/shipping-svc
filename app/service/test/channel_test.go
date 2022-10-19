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

var channelRepository = &repository_mock.ChannelRepositoryMock{Mock: mock.Mock{}}
var shippingCourierStatusRepository = &repository_mock.ShippingCourierStatusRepositoryMock{Mock: mock.Mock{}}
var channelSvc = service.NewChannelService(logger, baseRepository, channelRepository, shippingCourierStatusRepository)

// func init() {
// }

var channelTest = entity.Channel{
	ChannelName: "test",
	ChannelCode: "channel code test",
	Description: "description test",
	Status:      int32(1),
	Logo:        "logo test",
}

func TestCreateChannel(t *testing.T) {
	//status := int32(1)
	req := request.SaveChannelRequest{
		ChannelName: channelTest.ChannelName,
		ChannelCode: channelTest.ChannelCode,
		Description: channelTest.Description,
		Status:      &channelTest.Status,
		Logo:        channelTest.Logo,
	}
	channels := []entity.Channel{}
	channel := entity.Channel{}

	filter := map[string]interface{}{
		"channel_code": channelTest.ChannelCode,
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	channelRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(channels, &paginationResult)
	channelRepository.Mock.On("CreateChannel", &req).Return(channel)
	result, _ := channelSvc.CreateChannel(req)
	assert.NotNil(t, result)
	assert.Equal(t, channelTest.ChannelName, result.ChannelName, channelNameIsNotCorrect)
	assert.Equal(t, channelTest.ChannelCode, result.ChannelCode, channelCodeIsNotCorrect)
}

func TestGetChannel(t *testing.T) {
	channel := entity.Channel{
		ChannelCode: "string",
	}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	channelRepository.Mock.On("FindByUid", &uid).Return(channel)
	result, _ := channelSvc.GetChannel(uid)

	assert.NotNil(t, result, "Cannot nil")
	assert.Equal(t, "string", result.ChannelCode, "ChannelCode must be string")
}

func TestDeleteChannel(t *testing.T) {
	channel := entity.Channel{}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	channelRepository.Mock.On("FindByUid", &uid).Return(channel)
	channelRepository.Mock.On("IsChannelHasChild").Return(&entity.ChannelHasChildFlag{}).Once()
	msg := channelSvc.DeleteChannel(uid)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, codeIsNotCorrect)
}

func TestListChannel(t *testing.T) {
	req := request.ChannelListRequest{
		Page: 1,
		Sort: "",
		Filters: request.ChannelListFilter{
			ChannelCode: []string{},
		},
		Limit: 10,
	}

	channel := []entity.Channel{
		{
			ChannelCode: "string",
			ChannelName: "string1",
		},
	}

	filter := map[string]interface{}{
		"channel_code": req.Filters.ChannelCode,
		"channel_name": req.Filters.ChannelName,
		"status":       req.Filters.Status,
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	channelRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(channel, &paginationResult)
	channels, pagination, msg := channelSvc.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, codeIsNotCorrect)
	assert.Equal(t, 1, len(channels), "Count of Channels must be 1")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}

func TestCreateChannelFail(t *testing.T) {
	//status := int32(1)
	req := request.SaveChannelRequest{
		ChannelName: channelTest.ChannelName,
		ChannelCode: channelTest.ChannelCode,
		Description: channelTest.Description,
		Status:      &channelTest.Status,
		Logo:        channelTest.Logo,
	}
	channels := []entity.Channel{}
	channel := entity.Channel{}

	errTest := message.ErrDataChannelExists
	filter := map[string]interface{}{
		"channel_code": "string",
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	channelRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(channels, &paginationResult)
	channelRepository.Mock.On("CreateChannel", &req).Return(channel)
	channelSvc.CreateChannel(req)

	assert.Equal(t, message.ErrDataChannelExists.Code, errTest.Code, codeIsNotCorrect)
}

func TestUpdateChannelFail(t *testing.T) {
	req := request.UpdateChannelRequest{
		Uid:         "BnOI8D7p9rR7tI1R9rySw",
		ChannelName: channelTest.ChannelName,
		ChannelCode: channelTest.ChannelCode,
		Description: channelTest.Description,
		Status:      int(channelTest.Status),
		Logo:        channelTest.Logo,
	}
	channel := entity.Channel{}

	var isExist bool
	errTest := message.ErrDataChannelExists

	channelRepository.Mock.On("FindByUid", &req.Uid).Return(channel)
	channelRepository.Mock.On("CheckExistsByUIdChannelCode", req.Uid, req.ChannelCode).Return(isExist)
	channelRepository.Mock.On("UpdateChannel", &req).Return(channel)
	channelSvc.UpdateChannel(req)

	assert.Equal(t, message.ErrDataChannelExists.Code, errTest.Code, codeIsNotCorrect)
}
func TestGetChannelFail(t *testing.T) {
	channel := entity.Channel{}
	errTest := message.ErrNoData

	uid := "gj2MZ9CBfdfdhcHSNVOLpUeqUUUU"
	channelRepository.Mock.On("FindByUid", &uid).Return(channel, errTest)
	channelRepository.Mock.On("GetChannel", &uid).Return(channel)
	channelSvc.GetChannel(uid)

	errIsNotFound := "Data is not found"
	errCodeIsNotFound := 34005
	assert.EqualError(t, errors.New(errIsNotFound), errTest.Message, "Channel is not found")
	assert.Equal(t, errCodeIsNotFound, errTest.Code, "Channel is not found")
}

func TestGetListStatusSuccess(t *testing.T) {
	shippingCourierStatusRepository.Mock.On("FindByParams", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.ShippingCourierStatus{{}, {}}, &base.Pagination{}, nil).Once()
	input := request.GetChannelCourierStatusRequest{}
	result, _, msg := channelSvc.GetListStatus(input)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, codeIsNotCorrect)
	assert.Len(t, result, 2)
}

func TestGetListStatusNotFound(t *testing.T) {
	shippingCourierStatusRepository.Mock.On("FindByParams", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.ShippingCourierStatus{}, &base.Pagination{}, nil).Once()
	input := request.GetChannelCourierStatusRequest{}
	result, _, msg := channelSvc.GetListStatus(input)

	assert.Equal(t, message.SuccessMsg, msg)
	assert.Len(t, result, 0)
}

func TestDeleteChannelHasCourierAssigned(t *testing.T) {
	channel := entity.Channel{}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	channelRepository.Mock.On("FindByUid", &uid).Return(channel)
	channelRepository.Mock.On("IsChannelHasChild").Return(&entity.ChannelHasChildFlag{ChannelCourier: true}).Once()
	msg := channelSvc.DeleteChannel(uid)
	assert.Equal(t, message.ErrChannelHasCourierAssigned.Code, msg.Code, codeIsNotCorrect)
}

func TestDeleteChannelHasShippingStatus(t *testing.T) {
	channel := entity.Channel{}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	channelRepository.Mock.On("FindByUid", &uid).Return(channel)
	channelRepository.Mock.On("IsChannelHasChild").Return(&entity.ChannelHasChildFlag{ShippingStatus: true}).Once()
	msg := channelSvc.DeleteChannel(uid)
	assert.Equal(t, message.ErrChannelHasChildShippingStatus.Code, msg.Code, codeIsNotCorrect)
}

func TestDeleteChannelNotFound(t *testing.T) {
	uid := "BnOI8D7p9rR7tI1R9ryS"
	channelRepository.Mock.On("FindByUid", &uid).Return(nil)
	msg := channelSvc.DeleteChannel(uid)
	assert.Equal(t, message.ErrChannelNotFound.Code, msg.Code, codeIsNotCorrect)
}
