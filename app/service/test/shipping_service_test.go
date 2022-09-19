package test

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository/repository_mock"
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
var orderShippingRepository = &repository_mock.OrderShippingRepositoryMock{Mock: mock.Mock{}}

func init() {
	shippingService = service.NewShippingService(
		logger,
		baseRepository,
		channelRepository,
		courierServiceRepo,
		courierCoverageCodeRepository,
		shipper,
		redis,
		orderShippingRepository,
		courierRepository,
		shippingCourierStatusRepository,
	)
}

func TestGetShippingRate_ShipperSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
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
		CourierServiceUID: []string{"", ""},
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

func TestGetShippingRate_InternalSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
		TotalWeight:       10,
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal"},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant"},
		}).Once()

	courierCoverageCodeRepository.Mock.On("FindInternalAndMerchantCourierCoverage").Return(map[string]bool{"aa": true, "bb": true}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Len(t, result[0].Services, 2)
	assert.Equal(t, message.SuccessMsg.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.SuccessMsg.Message, result[0].Services[1].Error.Message)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_Internal_CoverageNotExistSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal"},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant"},
		}).Once()

	courierCoverageCodeRepository.Mock.On("FindInternalAndMerchantCourierCoverage").Return(map[string]bool{"aa": true}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, message.SuccessMsg.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.ErrCourierCoverageCodeUidNotExist.Message, result[0].Services[1].Error.Message)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_Internal_WieghtExceedSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
		TotalWeight:       10,
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal", MaxWeight: 1},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant"},
		}).Once()

	courierCoverageCodeRepository.Mock.On("FindInternalAndMerchantCourierCoverage").Return(map[string]bool{"aa": true}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrWeightExceeds.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.ErrCourierCoverageCodeUidNotExist.Message, result[0].Services[1].Error.Message)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_ShipperFailed(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{""},
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
		CourierServiceUID: []string{"", ""},
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
		CourierServiceUID: []string{"", ""},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrChannelNotFound, msg, codeIsNotCorrect)
}

func TestGetShippingRateCourierServiceUID_Required(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{},
	}

	result, msg := shippingService.GetShippingRate(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrCourierServiceIsRequired, msg, codeIsNotCorrect)
}

func TestGetShippingRateByShippingType_ShippingTypeRequired(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{},
	}

	result, msg := shippingService.GetShippingRateByShippingType(input)
	assert.Nil(t, result)
	assert.Equal(t, message.ErrShippingTypeRequired, msg, codeIsNotCorrect)
}

var createDeliveryRequest = &request.CreateDelivery{
	ChannelUID:        "chabcde",
	CouirerServiceUID: "csabcde",
	OrderNo:           "on12345",
	COD:               false,
	UseInsurance:      false,
	Merchant: request.CreateDeliveryPartner{
		Name:  "merchant name",
		UID:   "muid",
		Phone: "08980898",
		Email: "aaaa@aaa,com",
	},
	Customer: request.CreateDeliveryPartner{
		Name:  "customer name",
		UID:   "custuid",
		Phone: "08980877",
		Email: "bbb@bbb.com",
	},
	Origin: request.CreateDeiveryArea{
		Address:      "o add",
		CountryCode:  "ID",
		PostalCode:   "53355",
		Subdistrict:  "sdo",
		Latitude:     "",
		Longitude:    "",
		ProvinceCode: "pco",
		DistrictName: "do",
	},
	Destination: request.CreateDeiveryArea{
		Address:      "d add",
		CountryCode:  "ID",
		PostalCode:   "53356",
		Subdistrict:  "sdd",
		Latitude:     "",
		Longitude:    "",
		ProvinceCode: "pcd",
		DistrictName: "dd",
	},
	Notes: "",
	Package: request.CreateDeliveryPackage{
		Product: []request.CreateDeliveryProduct{
			{
				UID:   "aaauid",
				Name:  "aaan",
				Qty:   3,
				Price: 1000,
			},
			{
				UID:   "bbbuid",
				Name:  "bbbn",
				Qty:   5,
				Price: 2000,
			},
		},
		TotalWeight:         3,
		TotalWidth:          10,
		TotalLength:         20,
		TotalHeight:         30,
		TotalProductPrice:   13000,
		ContainPrescription: 0,
	},
}

func TestCreateDeliveryShipperSuccess(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus).Once()

	orderShipping := &entity.OrderShipping{
		BaseIDModel: base.BaseIDModel{
			ID:  5,
			UID: "osuid",
		},
		OrderNo: createDeliveryRequest.OrderNo,
	}

	order := &response.CreateDeliveryThirdPartyData{
		BookingID: "bookid",
		Status:    shipping_provider.StatusRequestPickup,
	}

	shipper.Mock.On("CreateDelivery", mock.Anything).Return(order, message.SuccessMsg).Once()
	orderShippingRepository.Mock.On("Upsert", mock.Anything).Return(orderShipping).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, createDeliveryRequest.OrderNo, result.OrderNoAPI)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCreateDeliveryShipperSaveFailed(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus).Once()

	orderShipping := &entity.OrderShipping{
		BaseIDModel: base.BaseIDModel{
			ID:  5,
			UID: "osuid",
		},
		OrderNo: createDeliveryRequest.OrderNo,
	}

	order := &response.CreateDeliveryThirdPartyData{
		BookingID: "bookid",
		Status:    shipping_provider.StatusRequestPickup,
	}

	shipper.Mock.On("CreateDelivery", mock.Anything).Return(order, message.SuccessMsg).Once()
	orderShippingRepository.Mock.On("Upsert", mock.Anything).Return(orderShipping, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestCreateDeliveryThridPartyCourierInvalid(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        "not shipper",
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierCode, msg)
}

func TestCreateDeliveryInvalidCourierType(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: "invalid",
		Code:        "not shipper",
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierType, msg)
}

func TestCreateDeliveryShipperFailed(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus).Once()

	order := &response.CreateDeliveryThirdPartyData{
		BookingID: "bookid",
		Status:    shipping_provider.StatusRequestPickup,
	}

	shipper.Mock.On("CreateDelivery", mock.Anything).Return(order, message.FailedMsg).Once()
	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestCreateDeliveryShipperShippingStatusNotFound(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	shippingStatus := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(shippingStatus, errors.New("")).Once()
	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrShippingStatus, msg)
}

func TestCreateDeliveryShipperOrderShippingAlreadyExist(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShipping := &entity.OrderShipping{
		BaseIDModel: base.BaseIDModel{
			ID:  5,
			UID: "osuid",
		},
		OrderNo: createDeliveryRequest.OrderNo,
		Status:  shipping_provider.StatusRequestPickup,
	}
	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(orderShipping).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderNoAlreadyExists, msg)
}

func TestCreateDeliveryShipperOrderShippingError(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrDB, msg)
}

func TestCreateDeliveryShipperCourierServiceNotFound(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(nil, nil).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCourierServiceNotFound, msg)
}

func TestCreateDeliveryShipperCourierServiceNotFoundError(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCourierServiceNotFound, msg)
}

func TestCreateDeliveryShipperChannelNotFound(t *testing.T) {
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrChannelNotFound, msg)
}

func TestCreateDeliveryShipperChannelNotFoundError(t *testing.T) {
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrChannelNotFound, msg)
}

var getOrderTrackingRequest = &request.GetOrderShippingTracking{
	UID:        "ORDER_SHIPPING_UID",
	ChannelUID: "CHHANNEL_UID",
}

func TestOrderShippingTrackingShipperSuccess(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: getOrderTrackingRequest.ChannelUID,
		},
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(&entity.OrderShipping{
			BaseIDModel: base.BaseIDModel{
				UID: getOrderTrackingRequest.UID,
			},
			CourierID: courier.ID,
			BookingID: "SHIPER_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	shipper.Mock.On("GetTracking", mock.Anything).
		Return([]response.GetOrderShippingTracking{}).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestOrderShippingTrackingShipperGetOrderDetailError(t *testing.T) {
	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: getOrderTrackingRequest.ChannelUID,
		},
	}

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(&entity.OrderShipping{
			BaseIDModel: base.BaseIDModel{
				UID: getOrderTrackingRequest.UID,
			},
			CourierID: courier.ID,
			BookingID: "SHIPER_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	shipper.Mock.On("GetTracking", mock.Anything).
		Return(nil, message.ErrGetOrderDetail).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrGetOrderDetail, msg)
}

func TestOrderShippingTrackingInvalidThirdPartyCourier(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        "INVALID",
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: getOrderTrackingRequest.ChannelUID,
		},
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(&entity.OrderShipping{
			BaseIDModel: base.BaseIDModel{
				UID: getOrderTrackingRequest.UID,
			},
			CourierID: courier.ID,
			BookingID: "SHIPER_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierCode, msg)
}

func TestOrderShippingTrackingInvalidCourierType(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: "INVALID",
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: getOrderTrackingRequest.ChannelUID,
		},
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(&entity.OrderShipping{
			BaseIDModel: base.BaseIDModel{
				UID: getOrderTrackingRequest.UID,
			},
			CourierID: courier.ID,
			BookingID: "SHIPER_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierType, msg)
}

func TestOrderShippingTrackingOrderNotBelongToChannel(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: "INVALID",
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "OTHER CHANNEL",
		},
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(&entity.OrderShipping{
			BaseIDModel: base.BaseIDModel{
				UID: getOrderTrackingRequest.UID,
			},
			CourierID: courier.ID,
			BookingID: "SHIPER_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderBelongToAnotherChannel, msg)
}

func TestOrderShippingTrackingGetOrderShippingNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestOrderShippingTrackingGetOrderShippingError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestOrderShippingTrackingChannelUIDRequired(t *testing.T) {
	result, msg := shippingService.OrderShippingTracking(&request.GetOrderShippingTracking{UID: "", ChannelUID: ""})
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrChannelUIDRequired, msg)
}

var updateStatusReq = &request.WebhookUpdateStatusShipper{
	ExternalID: "ORDERNO",
	ExternalStatus: request.ShipperStatus{
		Code:        1,
		Name:        "Name",
		Description: "Desc",
	},
	Awb: "AWB-00001",
}

func TestUpdateStatusShipper(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(&entity.OrderShipping{}).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateStatusShipperSaveFailed(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(nil, errors.New("")).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestUpdateStatusShipperShippingStatusNotFOund(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrShippingStatus, msg)
}

func TestUpdateStatusShipperGetShippingStatusError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil, errors.New("")).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrShippingStatus, msg)
}

func TestUpdateStatusShipperOrderNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(nil).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestUpdateStatusShipperGetOrderError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(nil, errors.New("")).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}
var getOrderShippingListRequest = request.GetOrderShippingList{
	Limit:   3,
	Page:    2,
	Sort:    "",
	Dir:     "desc",
	Filters: request.GetOrderShippingFilter{},
}

func TestGetOrderShippingList(t *testing.T) {
	orderShippingRepository.Mock.On("FindByParams", mock.Anything).
		Return([]response.GetOrderShippingList{},
			&base.Pagination{Limit: getOrderShippingListRequest.Limit, Page: getOrderShippingListRequest.Page}).
		Once()
	result, pagination, msg := shippingService.GetOrderShippingList(&getOrderShippingListRequest)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.NotNil(t, msg)
	assert.Equal(t, getOrderShippingListRequest.Limit, pagination.Limit)
	assert.Equal(t, getOrderShippingListRequest.Page, pagination.Page)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetOrderShippingListError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByParams", mock.Anything).
		Return(nil, nil, errors.New("")).
		Once()
	result, pagination, msg := shippingService.GetOrderShippingList(&getOrderShippingListRequest)
	assert.Nil(t, result)
	assert.Nil(t, pagination)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrNoData, msg)
}

func TestGetOrderShippingListInvalidDateFromError(t *testing.T) {
	req := getOrderShippingListRequest
	req.Filters.OrderShippingDateFrom = "INVALID_DATE"

	result, pagination, msg := shippingService.GetOrderShippingList(&req)
	assert.Nil(t, result)
	assert.Nil(t, pagination)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrFormatDateYYYYMMDD, msg)
}

func TestGetOrderShippingListInvalidDateToError(t *testing.T) {
	req := getOrderShippingListRequest
	req.Filters.OrderShippingDateTo = "INVALID_DATE"

	result, pagination, msg := shippingService.GetOrderShippingList(&req)
	assert.Nil(t, result)
	assert.Nil(t, pagination)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrFormatDateYYYYMMDD, msg)
}

func TestGetOrderShippingDetailByUIDSuccess(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID").Return(&entity.OrderShipping{
		Channel:              &entity.Channel{},
		Courier:              &entity.Courier{},
		CourierService:       &entity.CourierService{},
		OrderShippingItem:    []entity.OrderShippingItem{},
		OrderShippingHistory: []entity.OrderShippingHistory{},
	}).Once()

	shippingCourierStatusRepository.Mock.On("FindByCode").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()

	result, msg := shippingService.GetOrderShippingDetailByUID("")

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetOrderShippingDetailByUIDStatusNil(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID").Return(&entity.OrderShipping{
		Channel:              &entity.Channel{},
		Courier:              &entity.Courier{},
		CourierService:       &entity.CourierService{},
		OrderShippingItem:    []entity.OrderShippingItem{},
		OrderShippingHistory: []entity.OrderShippingHistory{},
	}).Once()

	shippingCourierStatusRepository.Mock.On("FindByCode").Return(nil).Once()

	result, msg := shippingService.GetOrderShippingDetailByUID("")

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetOrderShippingDetailByUIDNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID").Return(nil).Once()

	result, msg := shippingService.GetOrderShippingDetailByUID("")

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestGetOrderShippingDetailByUIDNotFoundError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID").Return(nil, errors.New("")).Once()

	result, msg := shippingService.GetOrderShippingDetailByUID("")

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}
