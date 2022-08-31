package test

import (
	"errors"
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
		Return(nil).Twice()

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "1"}).Times(4)

	shipper.Mock.On("GetShippingRate", mock.Anything).
		Return(&response.ShippingRateCommonResponse{
			Rate:    make(map[string]response.ShippingRateData),
			Summary: make(map[string]response.ShippingRateSummary),
		}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, msg, message.SuccessMsg, "Message is no correct")
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
			{CourierCode: ""},
			{CourierCode: ""},
		}).Once()

	redis.Mock.On("GetJsonStruct", mock.Anything).
		Return(nil).Twice()

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "1"}).Times(4)

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, msg, message.SuccessMsg, "Message is no correct")
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

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "1"}).Once()

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "2"}).Once()

	shipper.Mock.On("GetShippingRate", mock.Anything).
		Return(&response.ShippingRateCommonResponse{
			Rate:    make(map[string]response.ShippingRateData),
			Summary: make(map[string]response.ShippingRateSummary),
			Msg:     message.ErrGetShipperRate,
		}).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 400, result[0].Services[0].AvailableCode)
	assert.Equal(t, message.ErrGetShipperRate.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.SuccessMsg, msg, codeIsNotCorrect)
}

func TestGetShippingRate_DestinationNotFound(t *testing.T) {
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

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "1"}).Once()

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "2"}, errors.New("")).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 400, result[0].Services[0].AvailableCode)
	assert.Equal(t, message.ErrDestinationNotFound.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.SuccessMsg, msg, codeIsNotCorrect)
}

func TestGetShippingRate_OriginNotFound(t *testing.T) {
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

	courierCoverageCodeRepository.Mock.On("FindByCountryCodeAndPostalCode", mock.Anything).
		Return(&entity.CourierCoverageCode{Code1: "1"}, errors.New("")).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 400, result[0].Services[0].AvailableCode)
	assert.Equal(t, message.ErrOriginNotFound.Message, result[0].Services[0].Error.Message)
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
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, message.ErrCourierServiceNotFound, msg, "Message is no correct")
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
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, message.ErrChannelNotFound, msg, "Message is no correct")
}

func TestGetShippingRateCourierServiceUID_Required(t *testing.T) {
	input := request.GetShippingRateRequest{
		ChannelCourierService: []request.ChannelCourierServicePayloadItem{},
	}

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result, "Result should be nil")
	assert.Equal(t, message.ErrCourierServiceIsRequired, msg, "Message is no correct")
}
