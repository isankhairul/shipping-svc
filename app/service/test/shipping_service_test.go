package test

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/service"
	"testing"

	"go-klikdokter/helper/http_helper/shipping_provider"
	"go-klikdokter/helper/http_helper/shipping_provider/shipping_provider_mock"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/cache/cache_mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var shippingService service.ShippingService
var shipper = &shipping_provider_mock.ShipperMock{Mock: mock.Mock{}}
var redis = &cache_mock.Redis_Mock{Mock: mock.Mock{}}

func init() {
	shippingService = service.NewShippingService(
		logger,
		baseRepository,
		channelRepository,
		courierServiceRepo,
		courierCoverageCodeRepository,
		shipper,
		redis,
	)
}

func TestGetShippingRate_ShipperSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{
			{},
			{},
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierCode: shipping_provider.ShipperCode},
			{CourierCode: shipping_provider.ShipperCode},
		}).Once()

	redis.Mock.On("GetJsonStruct", mock.Anything).
		Return(nil).Once()

	shipper.Mock.On("GetShippingRate", mock.Anything).
		Return(&response.ShippingRateCommonResponse{
			Rate:    make(map[string]response.ShippingRateData),
			Summary: make(map[string]response.ShippingRateSummary),
		}).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_DefaultSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{
			{},
			{},
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierCode: "aa"},
			{CourierCode: "bb"},
		}).Once()

	redis.Mock.On("GetJsonStruct", mock.Anything).
		Return(nil).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_ShipperFailed(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{
			{},
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{{CourierCode: shipping_provider.ShipperCode}}).Once()

	redis.Mock.On("GetJsonStruct", mock.Anything).
		Return(nil).Once()

	shipper.Mock.On("GetShippingRate", mock.Anything).
		Return(&response.ShippingRateCommonResponse{
			Rate: make(map[string]response.ShippingRateData),
		}).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, 400, result[0].Services[0].AvailableCode)
	assert.Equal(t, message.ErrShippingRateNotFound.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.SuccessMsg, msg, codeIsNotCorrect)
}

func TestGetShippingRate_CourierServiceNotFound(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{
			{},
			{},
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{}).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrCourierServiceNotFound, msg, codeIsNotCorrect)
}

func TestGetShippingRate_ChannelNotFound(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{
			{},
			{},
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrChannelNotFound, msg, codeIsNotCorrect)
}

func TestGetShippingRateCourierServiceUID_Required(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{},
	}

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrCourierServiceIsRequired, msg, codeIsNotCorrect)
}
