package test

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var baseRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var courierRepository = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
var courierServiceRepository = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
var shipmentPredefinedRepository = &repository_mock.ShipmentPredefinedMock{Mock: mock.Mock{}}
var svc = service.NewCourierService(logger, baseRepository, courierRepository, courierServiceRepository, shipmentPredefinedRepository)

func init() {
}

func TestCreateCourierService(t *testing.T) {
	status := int32(1)
	req := request.SaveCourierServiceRequest{
		Cancelable:          1,
		CodAvailable:        1,
		CourierUId:          "gj2MZ9CBhcHSNVOLpUeqU",
		CreatedAt:           time.Now(),
		CreatedBy:           "Test",
		EndTime:             datatype.Time(""),
		ETD_Max:             1,
		ETD_Min:             1,
		Insurance:           1,
		InsuranceFee:        1,
		InsuranceFeeType:    "Test",
		InsuranceMin:        1,
		Logo:                "Test",
		MaxDistance:         1,
		MaxPurchase:         1,
		MaxVolume:           1,
		MaxWeight:           1,
		MinPurchase:         1,
		PrescriptionAllowed: 1,
		Repickup:            1,
		ShippingCode:        "Testing123456",
		ShippingDescription: "Test",
		ShippingName:        "Test",
		ShippingType:        "Test",
		StartTime:           datatype.Time(""),
		Status:              &status,
		TrackingAvailable:   1,
		UpdatedAt:           time.Now(),
		UpdatedBy:           "Test",
	}
	var isExist bool
	courier := entity.Courier{}
	courierService := entity.CourierService{}
	courierUId := req.CourierUId

	courierRepository.Mock.On("FindByUid", &courierUId).Return(courier)
	courierServiceRepository.Mock.On("CheckExistsByCourierIdShippingCode", req.CourierUId, req.ShippingCode).Return(isExist)
	courierServiceRepository.Mock.On("CreateCourierService", &req).Return(courierService)
	result, _ := svc.CreateCourierService(req)
	assert.NotNil(t, result)
	assert.Equal(t, "gj2MZ9CBhcHSNVOLpUeqU", result.CourierUId, "CourierUId must be gj2MZ9CBhcHSNVOLpUeqU")
	assert.Equal(t, "Testing123456", result.ShippingCode, "ShippingCode must be Testing123456")
}

func TestGetCourierService(t *testing.T) {
	CourierService := entity.CourierService{
		CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
		ShippingCode: "string",
	}

	uid := "gj2MZ9CBhcHSNVOLpUeqU"
	courierServiceRepository.Mock.On("FindByUid", &uid).Return(CourierService)
	result, _ := svc.GetCourierService(uid)

	assert.NotNil(t, result, "Cannot nil")
	assert.Equal(t, "gj2MZ9CBhcHSNVOLpUeqU", result.CourierUId, "CourierUId must be gj2MZ9CBhcHSNVOLpUeqU")
	assert.Equal(t, "string", result.ShippingCode, "ShippingCode must be string")
}

func TestDeleteCourierService(t *testing.T) {
	CourierService := entity.CourierService{
		BaseIDModel:  base.BaseIDModel{ID: 1},
		CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
		ShippingCode: "string",
	}

	uid := "gj2MZ9CBhcHSNVOLpUeqU"
	courierServiceRepository.Mock.On("FindByUid", &uid).Return(CourierService)
	courierServiceRepository.Mock.On("IsCourierServiceAssigned").Return(false).Once()
	msg := svc.DeleteCourierService(uid)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

func TestListCourierService(t *testing.T) {
	req := request.CourierServiceListRequest{
		Page:    1,
		Sort:    "",
		Filters: request.CourierServiceListFilter{CourierUID: []string{}},
		Limit:   10,
	}

	CourierService := []entity.CourierService{
		{
			CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
			ShippingCode: "string",
		},
		{
			CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
			ShippingCode: "string2",
		},
		{
			CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
			ShippingCode: "string3",
		},
	}

	filter := map[string]interface{}{
		"courier_uid":        req.Filters.CourierUID,
		"courier_type":       req.Filters.CourierType,
		"shipping_code":      req.Filters.ShippingCode,
		"shipping_name":      req.Filters.ShippingName,
		"shipping_type_code": req.Filters.ShippingTypeCode,
		"status":             req.Filters.Status,
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	courierServiceRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(CourierService, &paginationResult)
	CourierServices, pagination, msg := svc.GetListCourierService(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, 3, len(CourierServices), "Count of CourierServices must be 3")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}

func TestCreateCourierServiceFail(t *testing.T) {
	status := int32(1)
	req := request.SaveCourierServiceRequest{
		Cancelable:          1,
		CodAvailable:        1,
		CourierUId:          "gj2MZ9CBhcHSNVOLpUeqU",
		CreatedAt:           time.Now(),
		CreatedBy:           "Test",
		EndTime:             datatype.Time(""),
		ETD_Max:             1,
		ETD_Min:             1,
		Insurance:           1,
		InsuranceFee:        1,
		InsuranceFeeType:    "Test",
		InsuranceMin:        1,
		Logo:                "Test",
		MaxDistance:         1,
		MaxPurchase:         1,
		MaxVolume:           1,
		MaxWeight:           1,
		MinPurchase:         1,
		PrescriptionAllowed: 1,
		Repickup:            1,
		ShippingCode:        "string",
		ShippingDescription: "Test",
		ShippingName:        "Test",
		ShippingType:        "Test",
		StartTime:           datatype.Time(""),
		Status:              &status,
		TrackingAvailable:   1,
		UpdatedAt:           time.Now(),
		UpdatedBy:           "Test",
	}
	isExist := true
	courier := entity.Courier{}
	courierService := entity.CourierService{
		CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
		ShippingCode: "string",
	}
	courierUId := req.CourierUId

	courierRepository.Mock.On("FindByUid", &courierUId).Return(courier)
	courierServiceRepository.Mock.On("CheckExistsByCourierIdShippingCode", req.CourierUId, req.ShippingCode).Return(isExist)
	courierServiceRepository.Mock.On("CreateCourierService", &req).Return(courierService)
	_, err := svc.CreateCourierService(req)

	errIsExists := "Data courier_id/shipping_code already exists"
	errCodeIsExists := 34001
	assert.EqualError(t, errors.New(errIsExists), err.Message, "CourierUId and ShippingCode must be unique for each Courier")
	assert.Equal(t, errCodeIsExists, err.Code, "CourierUId and ShippingCode must be unique for each Courier")
}

func TestUpdateCourierServiceFail(t *testing.T) {
	req := request.UpdateCourierServiceRequest{
		Uid:                 "DYcO8MEsPJcuPIXlq30-T",
		Cancelable:          1,
		CodAvailable:        1,
		CourierUId:          "gj2MZ9CBhcHSNVOLpUeqU",
		EndTime:             datatype.Time(""),
		ETD_Max:             1,
		ETD_Min:             1,
		Insurance:           1,
		InsuranceFee:        1,
		InsuranceFeeType:    "Test",
		InsuranceMin:        1,
		Logo:                "Test",
		MaxDistance:         1,
		MaxPurchase:         1,
		MaxVolume:           1,
		MaxWeight:           1,
		MinPurchase:         1,
		PrescriptionAllowed: 1,
		Repickup:            1,
		ShippingCode:        "string2",
		ShippingDescription: "Test",
		ShippingName:        "Test",
		ShippingType:        "Test",
		StartTime:           datatype.Time(""),
		Status:              1,
		TrackingAvailable:   1,
	}
	var isExist bool
	courier := entity.Courier{}
	courierService := entity.CourierService{
		CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
		ShippingCode: "string2",
	}
	courierUId := req.CourierUId

	courierRepository.Mock.On("FindByUid", &courierUId).Return(courier)
	courierServiceRepository.Mock.On("FindByUid", &req.Uid).Return(courierService)
	courierServiceRepository.Mock.On("CheckExistsByCourierIdShippingCode", req.CourierUId, req.ShippingCode).Return(isExist)
	courierServiceRepository.Mock.On("UpdateCourierService", &req).Return(courierService)
	_, err := svc.UpdateCourierService(req.Uid, req)

	errIsExists := "Data courier_id/shipping_code already exists"
	errCodeIsExists := 34001
	assert.EqualError(t, errors.New(errIsExists), err.Message, "CourierUId and ShippingCode must be unique for each Courier")
	assert.Equal(t, errCodeIsExists, err.Code, "CourierUId and ShippingCode must be unique for each Courier")
}

func TestGetCourierServiceFail(t *testing.T) {
	CourierService := entity.CourierService{}
	errTest := message.ErrNoDataCourierService

	uid := "gj2MZ9CBfdfdhcHSNVOLpUeqUU"
	courierServiceRepository.Mock.On("FindByUid", &uid).Return(CourierService, errTest)
	courierServiceRepository.Mock.On("GetCourierService", &uid).Return(CourierService)
	svc.GetCourierService(uid)

	errIsNotFound := "Courier Service data not found"
	errCodeIsNotFound := 34005
	assert.EqualError(t, errors.New(errIsNotFound), errTest.Message, "Courier Service is not found")
	assert.Equal(t, errCodeIsNotFound, errTest.Code, "Courier Service is not found")
}

func TestDeleteCourierServiceHasAssigned(t *testing.T) {
	CourierService := entity.CourierService{
		CourierUId:   "gj2MZ9CBhcHSNVOLpUeqU",
		ShippingCode: "string",
	}

	uid := "gj2MZ9CBhcHSNVOLpUeqU"
	courierServiceRepository.Mock.On("FindByUid", &uid).Return(CourierService)
	courierServiceRepository.Mock.On("IsCourierServiceAssigned").Return(true).Once()
	msg := svc.DeleteCourierService(uid)

	assert.Equal(t, message.ErrCourierServiceHasAssigned.Code, msg.Code, "Code must be 201000")
}

func TestDeleteCourierServiceNotFound(t *testing.T) {
	uid := "gj2MZ9CBhcHSNVOLpUeq"
	courierServiceRepository.Mock.On("FindByUid", &uid).Return(nil)
	msg := svc.DeleteCourierService(uid)

	assert.Equal(t, message.ErrCourierServiceNotFound.Code, msg.Code, "Code must be 201000")
}

func TestGetCourierShippingTypeSuccess(t *testing.T) {
	shipmentPredefinedRepository.Mock.On("GetListByType").
		Return([]entity.ShippmentPredefined{
			{
				BaseIDModel: base.BaseIDModel{UID: "111"},
				Type:        "shiping_type",
			},
			{
				BaseIDModel: base.BaseIDModel{UID: "222"},
				Type:        "shiping_type",
			},
			{
				BaseIDModel: base.BaseIDModel{UID: "333"},
				Type:        "shiping_type",
			},
		}).Once()
	result, msg := svc.GetCourierShippingType()

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code is wrong")
	assert.Len(t, result, 3)
}

func TestGetCourierShippingTypeNotFound(t *testing.T) {
	shipmentPredefinedRepository.Mock.On("GetListByType").
		Return([]entity.ShippmentPredefined{}).Once()
	result, msg := svc.GetCourierShippingType()

	assert.Equal(t, message.ErrNoData.Code, msg.Code, "Code is wrong")
	assert.Len(t, result, 0)
}
