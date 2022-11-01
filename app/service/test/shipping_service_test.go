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

	"go-klikdokter/helper/http_helper/http_helper_mock"
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
var dapr = &http_helper_mock.DaprEndpointMock{Mock: mock.Mock{}}
var grab = &shipping_provider_mock.GrabMock{Mock: mock.Mock{}}

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
		dapr,
		grab,
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
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal",
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant",
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
		}).Once()

	courierCoverageCodeRepository.Mock.On("FindInternalAndMerchantCourierCoverage").Return(map[string]bool{"aa": true, "bb": true}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Len(t, result[0].Services, 2)
	assert.Equal(t, message.SuccessMsg.Message, result[0].Services[0].Error.Message)
	assert.Equal(t, message.SuccessMsg.Message, result[0].Services[1].Error.Message)
	assert.Equal(t, msg, message.SuccessMsg, codeIsNotCorrect)
}

func TestGetShippingRate_GrabSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
		Origin: request.AreaDetailPayload{
			Latitude:  "1",
			Longitude: "2",
		},
		Destination: request.AreaDetailPayload{
			Latitude:  "1",
			Longitude: "2",
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierCode: shipping_provider.GrabCode},
			{CourierCode: shipping_provider.GrabCode},
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

func TestGetShippingRate_GrabCoordinateRequired_Failed(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
		Origin: request.AreaDetailPayload{
			Latitude:  "",
			Longitude: "",
		},
		Destination: request.AreaDetailPayload{
			Latitude:  "",
			Longitude: "",
		},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierCode: shipping_provider.GrabCode},
			{CourierCode: shipping_provider.GrabCode},
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
	assert.Len(t, result, 1)
	assert.Equal(t, message.SuccessMsg, msg, codeIsNotCorrect)
	assert.Len(t, result[0].Services, 2)
}

func TestGetShippingRate_Internal_CoverageNotExistSuccess(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(entity.Channel{BaseIDModel: base.BaseIDModel{UID: "1"}}).Once()

	courierServiceRepo.Mock.On("FindCourierServiceByChannelAndUIDs", mock.Anything).
		Return([]entity.ChannelCourierServiceForShippingRate{
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal",
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant",
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
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
			{CourierID: 1, CourierCode: "aa", CourierTypeCode: "internal", MaxWeight: 1,
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
			{CourierID: 2, CourierCode: "bb", CourierTypeCode: "merchant",
				CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1},
		}).Once()

	courierCoverageCodeRepository.Mock.On("FindInternalAndMerchantCourierCoverage").Return(map[string]bool{"aa": true}).Twice()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, message.WeightExceedsMsg.Message, result[0].Services[0].Error.Message)
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
		Return([]entity.ChannelCourierServiceForShippingRate{{CourierCode: shipping_provider.ShipperCode,
			CourierStatus: 1, CourierServiceStatus: 1, ChannelCourierStatus: 1, ChannelCourierServiceStatus: 1, HidePurpose: 0, PrescriptionAllowed: 1}}).Once()

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
	assert.NotNil(t, result)
	assert.Equal(t, message.CourierServiceNotFoundMsg, msg, codeIsNotCorrect)
}

func TestGetShippingRate_ChannelNotFound(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{"", ""},
	}

	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrChannelNotFound, msg, codeIsNotCorrect)
}

func TestGetShippingRateCourierServiceUID_Required(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{},
	}

	result, msg := shippingService.GetShippingRate(input)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrCourierServiceIsRequired, msg, codeIsNotCorrect)
}

func TestGetShippingRateByShippingType_ShippingTypeRequired(t *testing.T) {
	input := request.GetShippingRateRequest{
		CourierServiceUID: []string{},
	}

	result, msg := shippingService.GetShippingRateByShippingType(input)
	assert.NotNil(t, result)
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
		ProvinceName: "pco",
		DistrictName: "do",
	},
	Destination: request.CreateDeiveryArea{
		Address:      "d add",
		CountryCode:  "ID",
		PostalCode:   "53356",
		Subdistrict:  "sdd",
		Latitude:     "",
		Longitude:    "",
		ProvinceName: "pcd",
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

var active int32 = 1

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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	pickup := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(pickup).Once()

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

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

func TestCreateDeliveryGrabSuccess(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.GrabCode,
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil).Once()

	pickup := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusRequestPickup,
		StatusCourier:    []byte(""),
	}

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(pickup).Once()

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

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

	grab.Mock.On("CreateDelivery", mock.Anything).Return(order, message.SuccessMsg).Once()
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

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

	assert.NotNil(t, result)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

	order := &response.CreateDeliveryThirdPartyData{
		BookingID: "bookid",
		Status:    shipping_provider.StatusRequestPickup,
	}

	shipper.Mock.On("CreateDelivery", mock.Anything).Return(order, message.FailedMsg).Once()
	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.FailedMsg, msg)
}

func TestCreateDeliveryGrabFailed(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.GrabCode,
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	created := &entity.ShippingCourierStatus{
		BaseIDModel: base.BaseIDModel{
			ID:  4,
			UID: "ssuid",
		},
		ShippingStatusID: 1,
		CourierID:        courier.ID,
		StatusCode:       shipping_provider.StatusCreated,
		StatusCourier:    []byte(""),
	}
	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(created).Once()

	order := &response.CreateDeliveryThirdPartyData{
		BookingID: "bookid",
		Status:    shipping_provider.StatusRequestPickup,
	}

	grab.Mock.On("CreateDelivery", mock.Anything).Return(order, message.FailedMsg).Once()
	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
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

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.OrderNoAlreadyExistsMsg, msg)
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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	orderShippingRepository.Mock.On("FindByOrderNo", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrDB, msg)
}

func TestCreateDeliveryShipperCoureirNotActive(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	var inactive int32 = 0
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: "cuid",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.ShipperCode,
		Status:      &inactive,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &active,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.CourierNotActiveMsg, msg)
}

func TestCreateDeliveryShipperCoureirServiceNotActive(t *testing.T) {

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
		Status:      &active,
	}
	var inactive int32 = 0
	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID: 2,
		Courier:   courier,
		Status:    &inactive,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.CourierServiceNotActiveMsg, msg)
}

func TestCreateDeliveryShipperCoureirServiceWeightExeeds(t *testing.T) {

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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID:           2,
		Courier:             courier,
		Status:              &active,
		MaxWeight:           1,
		PrescriptionAllowed: 1,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	req := *createDeliveryRequest
	req.Package.ContainPrescription = 1
	result, msg := shippingService.CreateDelivery(&req)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.WeightExceedsMsg, msg)
}

func TestCreateDeliveryShipperCoureirServicePrescriptionNotAllowed(t *testing.T) {

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
		Status:      &active,
	}

	courierService := &entity.CourierService{
		BaseIDModel: base.BaseIDModel{
			ID:  3,
			UID: createDeliveryRequest.CouirerServiceUID,
		},
		CourierID:           2,
		Courier:             courier,
		Status:              &active,
		PrescriptionAllowed: 0,
	}

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(courierService).Once()

	req := *createDeliveryRequest
	req.Package.ContainPrescription = 1
	result, msg := shippingService.CreateDelivery(&req)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.PrescriptionNotAllowedMsg, msg)
}

func TestCreateDeliveryShipperCourierServiceNotFound(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(nil, nil).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.CourierServiceNotFoundMsg, msg)
}

func TestCreateDeliveryShipperCourierServiceNotFoundError(t *testing.T) {

	channel := entity.Channel{BaseIDModel: base.BaseIDModel{ID: 1, UID: createDeliveryRequest.ChannelUID}}
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(channel).Once()

	courierServiceRepo.Mock.On("FindCourierService", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.CourierServiceNotFoundMsg, msg)
}

func TestCreateDeliveryShipperChannelNotFound(t *testing.T) {
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ChannelNotFoundMsg, msg)
}

func TestCreateDeliveryShipperChannelNotFoundError(t *testing.T) {
	channelRepository.Mock.On("FindByUid", mock.Anything).
		Return(nil, errors.New("")).Once()

	result, msg := shippingService.CreateDelivery(createDeliveryRequest)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ChannelNotFoundMsg, msg)
}

var getOrderTrackingRequest = &request.GetOrderShippingTracking{
	UID:        "ORDER_SHIPPING_UID",
	ChannelUID: "CHANNEL_UID",
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

func TestOrderShippingTrackingGrabSuccess(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.GrabCode,
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
			BookingID: "GRAB_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	grab.Mock.On("GetTracking", mock.Anything).
		Return([]response.GetOrderShippingTracking{}).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestOrderShippingTrackingGrabGetOrderDetailError(t *testing.T) {
	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{
			ID:  1,
			UID: "COURIER_UID",
		},
		CourierType: shipping_provider.ThirPartyCourier,
		Code:        shipping_provider.GrabCode,
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
			BookingID: "GRAB_ORDER_ID",
			Courier:   courier,
			Channel:   channel,
		}).Once()

	grab.Mock.On("GetTracking", mock.Anything).
		Return(nil, message.ErrGetOrderDetail).Once()

	result, msg := shippingService.OrderShippingTracking(getOrderTrackingRequest)
	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrGetOrderDetail, msg)
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
	assert.NotNil(t, result)
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
	Awb:  "AWB-00001",
	Auth: "466deec76ecdf5fca6d38571f6324d54",
}

func TestUpdateStatusShipper(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code: shipping_provider.ShipperCode,
		},
		CourierService: &entity.CourierService{},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(&entity.OrderShipping{
		Channel:        &entity.Channel{},
		Courier:        &entity.Courier{},
		CourierService: &entity.CourierService{},
	}).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.NotNil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateStatusShipperSaveFailed(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Courier: &entity.Courier{
			Code: shipping_provider.ShipperCode,
		},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(nil, errors.New("")).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestUpdateStatusShipperShippingStatusNotFOund(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Courier: &entity.Courier{
			Code: shipping_provider.ShipperCode,
		},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestUpdateStatusShipperGetShippingStatusError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Courier: &entity.Courier{
			Code: shipping_provider.ShipperCode,
		},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil, errors.New("")).Once()

	result, msg := shippingService.UpdateStatusShipper(updateStatusReq)

	assert.Nil(t, result)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
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
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrFormatDateYYYYMMDD, msg)
}

func TestGetOrderShippingListInvalidDateToError(t *testing.T) {
	req := getOrderShippingListRequest
	req.Filters.OrderShippingDateTo = "INVALID_DATE"

	result, pagination, msg := shippingService.GetOrderShippingList(&req)
	assert.NotNil(t, result)
	assert.NotNil(t, pagination)
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

var orderShipping = entity.OrderShipping{
	Channel: &entity.Channel{},
	Courier: &entity.Courier{
		Code:        shipping_provider.ShipperCode,
		CourierType: shipping_provider.ThirPartyCourier,
	},
	CourierService:       &entity.CourierService{Cancelable: 1},
	OrderShippingItem:    []entity.OrderShippingItem{},
	OrderShippingHistory: []entity.OrderShippingHistory{},
	Status:               shipping_provider.StatusRequestPickup,
	PickupCode:           new(string),
}

func TestCancelPickUpSuccess(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelPickupRequest", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(&order).Once()
	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCancelPickUpFailed(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelPickupRequest", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(nil, errors.New("")).Once()
	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestCancelPickUpShippingStatusNotFound(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelPickupRequest", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(nil).Once()
	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestCancelPickUpShippingThirdPartyError(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelPickupRequest", mock.Anything).Return(nil, errors.New("")).Once()

	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCancelPickup, msg)
}

func TestCancelPickUpShippingThirdPartyOrderNotCancelableError(t *testing.T) {
	order := orderShipping
	order.Status = "not_cancelable"
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCantCancelOrderShipping, msg)
}

func TestCancelPickUpShippingInvalidCourierTypeError(t *testing.T) {
	order := orderShipping
	order.Courier.CourierType = ""
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierType, msg)
}

func TestCancelPickUpShippingCourierServiceNotCancelableError(t *testing.T) {
	order := orderShipping
	order.CourierService.Cancelable = 0
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCantCancelOrderCourierService, msg)
}

func TestCancelPickUpShippingOrderServiceNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil).Once()
	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestCancelPickUpShippingOrderServiceNotFoundError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil, errors.New("")).Once()
	msg := shippingService.CancelPickup(&request.CancelPickup{UID: "uid"})
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

var cancelOrderReq = &request.CancelOrder{UID: "", Body: request.CancelOrderBodyRequest{Reason: "reason"}}

func TestCancelOrderSuccess(t *testing.T) {
	order := orderShipping
	order.CourierService.Cancelable = 1
	order.Courier.CourierType = shipping_provider.ThirPartyCourier
	order.Courier.Code = shipping_provider.ShipperCode

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelOrder", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(&order).Once()
	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestCancelOrderFailed(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelOrder", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(nil, errors.New("")).Once()
	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestCancelOrderShippingStatusNotFound(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelOrder", mock.Anything).Return(nil).Once()
	shippingCourierStatusRepository.Mock.On("FindByCode").Return(nil).Once()
	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestCancelOrderShippingThirdPartyError(t *testing.T) {
	order := orderShipping
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()
	shipper.Mock.On("CancelOrder", mock.Anything).Return(nil, errors.New("")).Once()

	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCancelPickup, msg)
}

func TestCancelOrderShippingThirdPartyOrderNotCancelableError(t *testing.T) {
	order := orderShipping
	order.Status = "not_cancelable"
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCantCancelOrderShipping, msg)
}

func TestCancelOrderShippingInvalidCourierTypeError(t *testing.T) {
	order := orderShipping
	order.Courier.CourierType = ""
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrInvalidCourierType, msg)
}

func TestCancelOrderShippingCourierServiceNotCancelableError(t *testing.T) {
	order := orderShipping
	order.CourierService.Cancelable = 0
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(&order).Once()

	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrCantCancelOrderCourierService, msg)
}

func TestCancelOrderShippingOrderServiceNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil).Once()
	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestCancelOrderShippingOrderServiceNotFoundError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil, errors.New("")).Once()
	msg := shippingService.CancelOrder(cancelOrderReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestGetOrderShippingLabelSuccess(t *testing.T) {
	req := request.GetOrderShippingLabel{
		ChannelUID: "",
		Body:       request.GetOrderShippingLabelBody{OrderShippingUID: []string{"", ""}},
	}
	orderShippingRepository.Mock.On("FindByUIDs", mock.Anything).Return([]entity.OrderShipping{
		{
			OrderShippingItem: []entity.OrderShippingItem{
				{}, {}, {},
			},
			Channel:        &entity.Channel{},
			Courier:        &entity.Courier{},
			CourierService: &entity.CourierService{},
		},
		{
			OrderShippingItem: []entity.OrderShippingItem{
				{},
			},
			Channel:        &entity.Channel{},
			Courier:        &entity.Courier{},
			CourierService: &entity.CourierService{},
		},
	}).Once()

	result, msg := shippingService.GetOrderShippingLabel(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].OrderShippingItems, 3)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetOrderShippingLabelHideItemsSuccess(t *testing.T) {
	req := request.GetOrderShippingLabel{
		ChannelUID: "",
		Body:       request.GetOrderShippingLabelBody{OrderShippingUID: []string{"", ""}, HideProduct: true},
	}
	orderShippingRepository.Mock.On("FindByUIDs", mock.Anything).Return([]entity.OrderShipping{
		{
			OrderShippingItem: []entity.OrderShippingItem{
				{}, {}, {},
			},
			Channel:        &entity.Channel{},
			Courier:        &entity.Courier{},
			CourierService: &entity.CourierService{},
		},
		{
			OrderShippingItem: []entity.OrderShippingItem{
				{},
			},
			Channel:        &entity.Channel{},
			Courier:        &entity.Courier{},
			CourierService: &entity.CourierService{},
		},
	}).Once()

	result, msg := shippingService.GetOrderShippingLabel(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Len(t, result[0].OrderShippingItems, 0)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestGetOrderShippingLabelError(t *testing.T) {
	req := request.GetOrderShippingLabel{
		ChannelUID: "",
		Body:       request.GetOrderShippingLabelBody{OrderShippingUID: []string{"", ""}, HideProduct: true},
	}
	orderShippingRepository.Mock.On("FindByUIDs", mock.Anything).Return(nil, errors.New("")).Once()

	result, msg := shippingService.GetOrderShippingLabel(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestRepickupSuccess(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        shipping_provider.ShipperCode,
			CourierType: shipping_provider.ThirPartyCourier,
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()
	shipper.Mock.On("CreatePickUpOrderWithTimeSlots", mock.Anything).
		Return(&response.CreatePickUpOrderShipperResponse{
			Data: response.CreatePickUpOrderShipper{
				OrderActivation: []response.CreatePickUpOrderOrderActivation{
					{PickUpCode: "001"},
				},
			},
		}).Once()

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(&entity.ShippingCourierStatus{}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(ordershipping).Once()
	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestRepickupSaveFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        shipping_provider.ShipperCode,
			CourierType: shipping_provider.ThirPartyCourier,
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()
	shipper.Mock.On("CreatePickUpOrderWithTimeSlots", mock.Anything).
		Return(&response.CreatePickUpOrderShipperResponse{
			Data: response.CreatePickUpOrderShipper{
				OrderActivation: []response.CreatePickUpOrderOrderActivation{
					{PickUpCode: "001"},
				},
			},
		}).Once()

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(&entity.ShippingCourierStatus{}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(ordershipping, errors.New("")).Once()
	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestRepickupShippingStatusNotFoundFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        shipping_provider.ShipperCode,
			CourierType: shipping_provider.ThirPartyCourier,
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()
	shipper.Mock.On("CreatePickUpOrderWithTimeSlots", mock.Anything).
		Return(&response.CreatePickUpOrderShipperResponse{
			Data: response.CreatePickUpOrderShipper{
				OrderActivation: []response.CreatePickUpOrderOrderActivation{
					{PickUpCode: "001"},
				},
			},
		}).Once()

	shippingCourierStatusRepository.Mock.On("FindByCode", mock.Anything).Return(nil).Once()
	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestRepickupShippingErroRequestShipperFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        shipping_provider.ShipperCode,
			CourierType: shipping_provider.ThirPartyCourier,
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()
	shipper.Mock.On("CreatePickUpOrderWithTimeSlots", mock.Anything).
		Return(nil, message.ErrCreatePickUpOrder).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrCreatePickUpOrder, msg)
}

func TestRepickupShippingInvalidCourierFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        "",
			CourierType: shipping_provider.ThirPartyCourier,
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrInvalidCourierCode, msg)
}

func TestRepickupShippingInvalidCourierTypeFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        "",
			CourierType: "",
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrInvalidCourierType, msg)
}

func TestRepickupShippingOrderHasBeenCancelledFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        "",
			CourierType: "",
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCancelled,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.OrderHasBeenCancelledMsg, msg)
}

func TestRepickupShippingRequestPickupHasBeenMadeFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code:        "",
			CourierType: "",
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusRequestPickup,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.RequestPickupHasBeenMadeMsg, msg)
}

func TestRepickupShippingOrderBelongsToAnotherChannelFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	ordershipping := &entity.OrderShipping{
		OrderShippingItem: []entity.OrderShippingItem{
			{}, {}, {},
		},
		Channel: &entity.Channel{BaseIDModel: base.BaseIDModel{
			UID: "A",
		}},
		Courier: &entity.Courier{
			Code:        "",
			CourierType: "",
		},
		CourierService: &entity.CourierService{},
		Status:         shipping_provider.StatusCreated,
		PickupCode:     new(string),
	}
	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(ordershipping).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrOrderBelongToAnotherChannel, msg)
}

func TestRepickupShippingOrderNotFoundFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestRepickupShippingErrorgetShippingOrderFailed(t *testing.T) {
	req := request.RepickupOrderRequest{
		ChannelUID:       "",
		OrderShippingUID: "",
		Username:         "",
	}

	orderShippingRepository.Mock.On("FindByUID", mock.Anything).Return(nil, errors.New("")).Once()

	result, msg := shippingService.RepickupOrder(&req)
	assert.NotNil(t, msg)
	assert.NotNil(t, result)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

var updateStatusGrabReq = &request.WebhookUpdateStatusGrabRequest{}

func TestUpdateStatusGrab(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code: shipping_provider.GrabCode,
		},
		CourierService: &entity.CourierService{},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(&entity.OrderShipping{
		Channel:        &entity.Channel{},
		Courier:        &entity.Courier{},
		CourierService: &entity.CourierService{},
	}).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)

	assert.NotNil(t, msg)
	assert.Equal(t, message.SuccessMsg, msg)
}

func TestUpdateStatusGrabSaveFailed(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code: shipping_provider.GrabCode,
		},
		CourierService: &entity.CourierService{},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(&entity.ShippingCourierStatus{
		ShippingStatus: &entity.ShippingStatus{},
	}).Once()
	orderShippingRepository.Mock.On("Upsert").Return(nil, errors.New("")).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrSaveOrderShipping, msg)
}

func TestUpdateStatusGrabShippingStatusNotFOund(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code: shipping_provider.GrabCode,
		},
		CourierService: &entity.CourierService{},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestUpdateStatusGrabGetShippingStatusError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(&entity.OrderShipping{
		Channel: &entity.Channel{},
		Courier: &entity.Courier{
			Code: shipping_provider.GrabCode,
		},
		CourierService: &entity.CourierService{},
	}).Once()
	shippingCourierStatusRepository.Mock.On("FindByCourierStatus").Return(nil, errors.New("")).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ShippingStatusNotFoundMsg, msg)
}

func TestUpdateStatusGrobOrderNotFound(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(nil).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}

func TestUpdateStatusGrabGetOrderError(t *testing.T) {
	orderShippingRepository.Mock.On("FindByOrderNo").Return(nil, errors.New("")).Once()

	msg := shippingService.UpdateStatusGrab(updateStatusGrabReq)
	assert.NotNil(t, msg)
	assert.Equal(t, message.ErrOrderShippingNotFound, msg)
}
