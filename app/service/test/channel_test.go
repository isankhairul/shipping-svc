package test

import (
	"errors"
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

var channelRepository = &repository_mock.ChannelRepositoryMock{Mock: mock.Mock{}}
var channelSvc = service.NewChannelService(logger, baseRepository, channelRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	//db.AutoMigrate(&entity.Channel{})
}

func TestCreateChannel(t *testing.T) {
	req := request.SaveChannelRequest{
		ChannelName: "test",
		ChannelCode: "channel code test",
		Description: "description test",
		Status:      1,
		Logo:        "logo test",
	}
	channels := []entity.Channel{}
	channel := entity.Channel{}

	filter := map[string]interface{}{
		"channel_code": "channel code test",
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
	assert.Equal(t, "test", result.ChannelName, "ChannelName must be test")
	assert.Equal(t, "channel code test", result.ChannelCode, "ChannelName must be channel code test")
	assert.Equal(t, "description test", result.Description, "Description must be description test")
	assert.Equal(t, 1, result.Status, "Status must be 1")
	assert.Equal(t, "logo test", result.Logo, "Status must be logo test")

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
	msg := channelSvc.DeleteChannel(uid)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

func TestListChannel(t *testing.T) {
	req := request.ChannelListRequest{
		Page:        1,
		Sort:        "",
		ChannelCode: "",
		Limit:       10,
	}

	channel := []entity.Channel{
		{
			ChannelCode: "string",
			ChannelName: "string1",
		},
	}

	filter := map[string]interface{}{
		"channel_code": "",
		"channel_name": "",
		"status":       0,
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	channelRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(channel, &paginationResult)
	channels, pagination, msg := channelSvc.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, 1, len(channels), "Count of Channels must be 1")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}

func TestCreateChannelFail(t *testing.T) {
	req := request.SaveChannelRequest{
		ChannelName: "test",
		ChannelCode: "string",
		Description: "description test",
		Status:      1,
		Logo:        "logo test",
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

	errIsExists := "Data channel_code already exists"
	errCodeIsExists := 34001
	assert.EqualError(t, errors.New(errIsExists), errTest.Message, "Channel Code must be unique")
	assert.Equal(t, errCodeIsExists, errTest.Code, "Channel Code must be unique")
}

func TestUpdateChannelFail(t *testing.T) {
	req := request.UpdateChannelRequest{
		Uid:         "BnOI8D7p9rR7tI1R9rySw",
		ChannelName: "test",
		ChannelCode: "string1",
		Description: "description test",
		Status:      1,
		Logo:        "logo test",
	}
	channel := entity.Channel{}

	var isExist bool
	errTest := message.ErrDataChannelExists

	channelRepository.Mock.On("FindByUid", &req.Uid).Return(channel)
	channelRepository.Mock.On("CheckExistsByUIdChannelCode", req.Uid, req.ChannelCode).Return(isExist)
	channelRepository.Mock.On("UpdateChannel", &req).Return(channel)
	channelSvc.UpdateChannel(req)

	errIsExists := "Data channel_code already exists"
	errCodeIsExists := 34001
	assert.EqualError(t, errors.New(errIsExists), errTest.Message, "Channel Code must be unique")
	assert.Equal(t, errCodeIsExists, errTest.Code, "Channel Code must be unique")
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
