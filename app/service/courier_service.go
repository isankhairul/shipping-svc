package service

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/gorm"
)

type CourierService interface {
	CreateCourier(input request.SaveCourierRequest) (*entity.Courier, message.Message)
	GetList(input request.CourierListRequest) ([]response.CourierListResponse, *base.Pagination, message.Message)
	UpdateCourier(uid string, input request.UpdateCourierRequest) (*entity.Courier, message.Message)
	GetCourier(uid string) (*entity.Courier, message.Message)
	DeleteCourier(uid string) message.Message

	//Courier-Service
	CreateCourierService(input request.SaveCourierServiceRequest) (*entity.CourierService, message.Message)
	//GetListCourierService(input request.CourierServiceListRequest) ([]*entity.CourierServiceDetailDTO, *base.Pagination, message.Message)
	GetListCourierService(input request.CourierServiceListRequest) ([]response.CourierServiceListResponse, *base.Pagination, message.Message)
	UpdateCourierService(uid string, input request.UpdateCourierServiceRequest) (*entity.CourierService, message.Message)
	GetCourierService(uid string) (*entity.CourierServiceDetailDTO, message.Message)
	DeleteCourierService(uid string) message.Message
	GetCourierShippingType() ([]response.ShippingTypeItem, message.Message)
}

type courierServiceImpl struct {
	logger                 log.Logger
	baseRepo               repository.BaseRepository
	courierRepo            repository.CourierRepository
	courierServiceRepo     repository.CourierServiceRepository
	shipmentPredefinedRepo repository.ShipmentPredefinedRepository
}

func NewCourierService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.CourierRepository,
	pcrp repository.CourierServiceRepository,
	sprp repository.ShipmentPredefinedRepository,
) CourierService {
	return &courierServiceImpl{lg, br, pr, pcrp, sprp}
}

// swagger:route POST /courier/courier Couriers SaveCourierRequest
// Create a new Courier
//
// responses:
//  401: errorResponse
//  500: errorResponse
//  201: Courier
func (s *courierServiceImpl) CreateCourier(input request.SaveCourierRequest) (*entity.Courier, message.Message) {
	logger := log.With(s.logger, "CourierService", "CreateCourier")
	courier, err := s.courierRepo.FindByCode(input.Code)

	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if courier != nil {
		return nil, message.ErrDuplicatedCourier
	}

	//Set request to entity
	Courier := entity.Courier{
		CourierName:           input.CourierName,
		Code:                  input.Code,
		CourierType:           input.CourierType,
		Logo:                  input.Logo,
		ProvideAirwaybill:     input.ProvideAirwaybill,
		CourierApiIntegration: input.CourierApiIntegration,
		HidePurpose:           input.HidePurpose,
		UseGeocoodinate:       input.UseGeocoodinate,
		Status:                &input.Status,
		ImageUID:              input.ImageUID,
		ImagePath:             input.ImagePath,
	}

	result, err := s.courierRepo.CreateCourier(&Courier)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	return result, message.SuccessMsg
}

// swagger:route GET /courier/courier/{uid} Couriers CourierByUIdParam
// Get Courier
//
// responses:
//  401: SuccessResponse
//  200: Courier
func (s *courierServiceImpl) GetCourier(uid string) (*entity.Courier, message.Message) {
	logger := log.With(s.logger, "CourierService", "GetCourier")

	result, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrCourierNotFound
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return result, message.SuccessMsg
}

// swagger:route GET /courier/courier Couriers CourierListRequest
// Get list of couriers
//
// responses:
//  401: errorResponse
//  200: SuccessResponse
func (s *courierServiceImpl) GetList(input request.CourierListRequest) ([]response.CourierListResponse, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	filter := map[string]interface{}{
		"status":       input.Filters.Status,
		"courier_type": input.Filters.CourierType,
		"code":         input.Filters.CourierCode,
		"courier_name": input.Filters.CourierName,
	}

	result, pagination, err := s.courierRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, nil, message.FailedMsg
	}

	if result == nil {
		_ = level.Warn(logger).Log(message.ErrNoData)
		return nil, nil, message.FailedMsg
	}

	return result, pagination, message.SuccessMsg
}

// swagger:route PUT /courier/courier/{uid} Couriers UpdateCourierRequest
// Update courier with specifeied id
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceImpl) UpdateCourier(uid string, input request.UpdateCourierRequest) (*entity.Courier, message.Message) {
	logger := log.With(s.logger, "CourierService", "UpdateCourier")

	_, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrCourierNotFound
	}

	data := map[string]interface{}{
		"courier_name":       input.CourierName,
		"courier_type":       input.CourierType,
		"code":               input.Code,
		"logo":               input.Logo,
		"status":             input.Status,
		"use_geocoodinate":   input.UseGeocoodinate,
		"provide_airwaybill": input.ProvideAirwaybill,
		"api_integration":    input.CourierApiIntegration,
		"hide_purpose":       input.HidePurpose,
		"image_uid":          input.ImageUID,
		"image_path":         input.ImagePath,
	}

	err = s.courierRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrCourierNotFound
	}
	return s.GetCourier(uid)
}

// swagger:route DELETE /courier/courier/{uid} Couriers DeleteCourierByUIdParam
// Delete courier by uid
//
// responses:
//  200: SuccessResponse
//  400: errorResponse
//  500: InternalServerErrorResponse
func (s *courierServiceImpl) DeleteCourier(uid string) message.Message {
	logger := log.With(s.logger, "CourierService", "DeleteCourier")

	courier, _ := s.courierRepo.FindByUid(&uid)

	if courier == nil {
		return message.ErrCourierNotFound
	}

	hasChild := s.courierRepo.IsCourierHasChild(courier.ID)

	if hasChild.CourierService {
		return message.ErrCourierHasChildCourierService
	}

	if hasChild.CourierCoverageCode {
		return message.ErrCourierHasChildCourierCoverage
	}

	if hasChild.ChannelCourier {
		return message.ErrCourierHasChildChannelCourier
	}

	if hasChild.ShippingCourierStatus {
		return message.ErrCourierHasChildShippingStatus
	}

	err := s.courierRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.ErrDB
	}

	return message.SuccessMsg
}

// swagger:route POST /courier/courier-services Courier-Services SaveCourierServiceRequest
// Add Courier Services
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) CreateCourierService(input request.SaveCourierServiceRequest) (*entity.CourierService, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "CreateCourierService")
	//Check exist courier_uid update
	courier, err := s.courierRepo.FindByUid(&input.CourierUId)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, message.ErrDataCourierUIdNotExist
		}
		return nil, message.FailedMsg
	}
	if courier == nil {
		return nil, message.ErrDataCourierUIdNotExist
	}
	//Check exists duplicate courier_uid/shipping_code
	isExists, err := s.courierServiceRepo.CheckExistsByCourierIdShippingCode(input.CourierUId, input.ShippingCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}
	if isExists {
		return nil, message.ErrDataCourierServiceExists
	}

	//Set request to entity
	defaultStatus := int32(1)
	courierService := entity.CourierService{
		//General
		CourierID:           courier.ID,
		CourierUId:          input.CourierUId,
		ShippingCode:        input.ShippingCode,
		ShippingName:        input.ShippingName,
		ShippingType:        input.ShippingType,
		ShippingDescription: input.ShippingDescription,
		ETD_Min:             input.ETD_Min,
		ETD_Max:             input.ETD_Max,
		Logo:                input.Logo,
		CodAvailable:        input.CodAvailable,
		PrescriptionAllowed: input.PrescriptionAllowed,
		Cancelable:          input.Cancelable,
		TrackingAvailable:   input.TrackingAvailable,
		Status:              &defaultStatus, //Default

		//Miscellaneous
		MaxWeight:        input.MaxWeight,
		MaxVolume:        input.MaxVolume,
		MaxDistance:      input.MaxDistance,
		MinPurchase:      input.MinPurchase,
		MaxPurchase:      input.MaxPurchase,
		Insurance:        input.Insurance,
		InsuranceMin:     input.InsuranceMin,
		InsuranceFeeType: input.InsuranceFeeType,
		InsuranceFee:     input.InsuranceFee,
		StartTime:        input.StartTime,
		EndTime:          input.EndTime,
		Repickup:         input.Repickup,
		ImageUID:         input.ImageUID,
		ImagePath:        input.ImagePath,
	}
	if input.Status != nil {
		courierService.Status = input.Status
	}
	resultInsert, err := s.courierServiceRepo.CreateCourierService(&courierService)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	return resultInsert, message.SuccessMsg
}

// swagger:route GET /courier/courier-services/{uid} Courier-Services CourierServiceRequestGetByUid
// Detail Courier Services
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) GetCourierService(uid string) (*entity.CourierServiceDetailDTO, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "GetCourierService")

	result, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, message.ErrNoDataCourierService
		}
		return nil, message.ErrDB
	}
	if result == nil {
		return nil, message.ErrNoDataCourierService
	}

	return ToCourierServiceDetailDTO(result), message.SuccessMsg
}

// swagger:route GET /courier/courier-services Courier-Services CourierServiceListRequest
// List of Courier Services
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) GetListCourierService(input request.CourierServiceListRequest) ([]response.CourierServiceListResponse, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	filter := map[string]interface{}{
		"courier_uid":   input.Filters.CourierUID,
		"courier_type":  input.Filters.CourierType,
		"shipping_code": input.Filters.ShippingCode,
		"shipping_name": input.Filters.ShippingName,
		"shipping_type": input.Filters.ShippingTypeCode,
		"status":        input.Filters.Status,
	}
	result, pagination, err := s.courierServiceRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, nil, message.FailedMsg
	}

	if result == nil {
		_ = level.Warn(logger).Log(message.ErrNoData)
		return nil, nil, message.ErrNoDataCourierService
	}

	//return convertToDTO(result), pagination, message.SuccessMsg
	return result, pagination, message.SuccessMsg
}

// swagger:route PUT /courier/courier-services/{uid} Courier-Services UpdateCourierServiceRequest
// Update Courier Services
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceImpl) UpdateCourierService(uid string, input request.UpdateCourierServiceRequest) (*entity.CourierService, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "UpdateCourierService")
	//Check exist courierServiceUId
	courierService, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, message.ErrDataCourierServiceUidNotExist
		}
		return nil, message.FailedMsg
	}
	if courierService == nil {
		return nil, message.ErrDataCourierServiceUidNotExist
	}
	//Check exist courier_uid update
	courier, err := s.courierRepo.FindByUid(&input.CourierUId)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, message.ErrDataCourierUIdNotExist
		}
		return nil, message.FailedMsg
	}
	if courier == nil {
		return nil, message.ErrDataCourierUIdNotExist
	}

	//Check exists duplicate courier_uid/shipping_code
	isExists, err := s.courierServiceRepo.CheckExistsByUIdCourierIdShippingCode(uid, input.CourierUId, input.ShippingCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}
	if isExists {
		return nil, message.ErrDataCourierServiceExists
	}

	data := map[string]interface{}{
		"courier_id":           courier.ID,
		"courier_uid":          input.CourierUId,
		"shipping_code":        input.ShippingCode,
		"shipping_name":        input.ShippingName,
		"shipping_type":        input.ShippingType,
		"shipping_description": input.ShippingDescription,
		"etd_min":              input.ETD_Min,
		"etd_max":              input.ETD_Max,
		"logo":                 input.Logo,
		"cod_available":        input.CodAvailable,
		"prescription_allowed": input.PrescriptionAllowed,
		"cancelable":           input.Cancelable,
		"tracking_available":   input.TrackingAvailable,
		"status":               input.Status,
		"max_weight":           input.MaxWeight,
		"max_volume":           input.MaxVolume,
		"max_distance":         input.MaxDistance,
		"min_purchase":         input.MinPurchase,
		"max_purchase":         input.MaxPurchase,
		"insurance":            input.Insurance,
		"insurance_min":        input.InsuranceMin,
		"insurance_fee_type":   input.InsuranceFeeType,
		"insurance_fee":        input.InsuranceFee,
		"start_time":           input.StartTime,
		"end_time":             input.EndTime,
		"repickup":             input.Repickup,
		"image_uid":            input.ImageUID,
		"image_path":           input.ImagePath,
	}

	err = s.courierServiceRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}

	//Find uid and response
	result, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, message.ErrNoDataCourierService
		}
		return nil, message.ErrDB
	}
	if result == nil {
		return nil, message.ErrNoDataCourierService
	}
	return result, message.SuccessMsg
}

// swagger:route DELETE /courier/courier-services/{uid} Courier-Services CourierServiceRequestDeleteByUid
// Delete Courier Services
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceImpl) DeleteCourierService(uid string) message.Message {
	logger := log.With(s.logger, "CourierServiceService", "DeleteCourierService")
	courierService, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return message.ErrNoDataCourierService
		}
		return message.FailedMsg
	}

	if courierService == nil {
		return message.ErrCourierServiceNotFound
	}

	if isAssigned := s.courierServiceRepo.IsCourierServiceAssigned(courierService.ID); isAssigned {
		return message.ErrCourierServiceHasAssigned
	}

	err = s.courierServiceRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}

// swagger:route GET /courier/shipping-type Couriers GetCourierShippingType
// Get List of Shipping Type
//
// responses:
//  200: ShippingTypeList
func (s *courierServiceImpl) GetCourierShippingType() ([]response.ShippingTypeItem, message.Message) {
	logger := log.With(s.logger, "CourierService", "GetCourierShippingType")

	result, err := s.shipmentPredefinedRepo.GetListByType("shipping_type")

	if err != nil {
		_ = level.Error(logger).Log("error", err.Error())
		return nil, message.ErrDB
	}

	if len(result) == 0 {
		return nil, message.ErrNoData
	}

	return response.NewShippingTypeItemList(result), message.SuccessMsg
}

func convertToDTO(services []entity.CourierService) []*entity.CourierServiceDetailDTO {
	items := make([]*entity.CourierServiceDetailDTO, len(services))
	for index, value := range services {
		items[index] = ToCourierServiceDetailDTO(&value)
	}
	return items
}

func ToCourierServiceDetailDTO(cs *entity.CourierService) *entity.CourierServiceDetailDTO {
	ret := &entity.CourierServiceDetailDTO{
		Uid:                 cs.UID,
		CourierUId:          cs.CourierUId,
		ShippingCode:        cs.ShippingCode,
		ShippingName:        cs.ShippingName,
		ShippingType:        cs.ShippingType,
		ShippingDescription: cs.ShippingDescription,
		ETD_Min:             cs.ETD_Min,
		ETD_Max:             cs.ETD_Max,
		Logo:                cs.Logo,
		CodAvailable:        cs.CodAvailable,
		PrescriptionAllowed: cs.PrescriptionAllowed,
		Cancelable:          cs.Cancelable,
		TrackingAvailable:   cs.TrackingAvailable,
		Status:              cs.Status,
		MaxWeight:           cs.MaxWeight,
		MaxVolume:           cs.MaxVolume,
		MaxDistance:         cs.MaxDistance,
		MinPurchase:         cs.MinPurchase,
		MaxPurchase:         cs.MaxPurchase,
		Insurance:           cs.Insurance,
		InsuranceMin:        cs.InsuranceMin,
		InsuranceFeeType:    cs.InsuranceFeeType,
		InsuranceFee:        cs.InsuranceFee,
		StartTime:           cs.StartTime,
		EndTime:             cs.EndTime,
		Repickup:            cs.Repickup,
		ImageUID:            cs.ImageUID,
		ImagePath:           cs.ImagePath,
	}
	if cs.Courier != nil {
		ret.CourierName = cs.Courier.CourierName
		ret.CourierType = cs.Courier.CourierType
		ret.CourierUId = cs.Courier.UID
	}

	return ret
}
