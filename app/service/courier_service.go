package service

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/gorm"
)

type CourierService interface {
	CreateCourier(input request.SaveCourierRequest) (*entity.Courier, message.Message)
	GetList(input request.CourierListRequest) ([]entity.Courier, *base.Pagination, message.Message)
	UpdateCourier(uid string, input request.UpdateCourierRequest) message.Message
	GetCourier(uid string) (*entity.Courier, message.Message)
	DeleteCourier(uid string) message.Message

	//Courier-Service
	CreateCourierService(input request.SaveCourierServiceRequest) (*entity.CourierService, message.Message)
	GetListCourierService(input request.CourierServiceListRequest) ([]entity.CourierService, *base.Pagination, message.Message)
	UpdateCourierService(uid string, input request.UpdateCourierServiceRequest) (*entity.CourierService, message.Message)
	GetCourierService(uid string) (*entity.CourierService, message.Message)
	DeleteCourierService(uid string) message.Message
}

type courierServiceImpl struct {
	logger             log.Logger
	baseRepo           repository.BaseRepository
	courierRepo        repository.CourierRepository
	courierServiceRepo repository.CourierServiceRepository
}

func NewCourierService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.CourierRepository,
	pcrp repository.CourierServiceRepository,
) CourierService {
	return &courierServiceImpl{lg, br, pr, pcrp}
}

// swagger:route POST /courier/courier Courier SaveCourierRequest
// Manage Courier
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) CreateCourier(input request.SaveCourierRequest) (*entity.Courier, message.Message) {
	logger := log.With(s.logger, "CourierService", "CreateCourier")
	s.baseRepo.BeginTx()
	//Set request to entity
	Courier := entity.Courier{
		CourierName: input.CourierName,
		Code:        input.Code,
		CourierType: input.CourierType,
		Logo:        input.Logo,
	}

	result, err := s.courierRepo.CreateCourier(&Courier)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.ErrDB
	}
	s.baseRepo.CommitTx()

	return result, message.SuccessMsg
}

// swagger:route GET /courier/courier/{uid} Get-Courier Courier
// Get Courier
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) GetCourier(uid string) (*entity.Courier, message.Message) {
	logger := log.With(s.logger, "CourierService", "GetCourier")

	result, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return result, message.SuccessMsg
}

func (s *courierServiceImpl) GetList(input request.CourierListRequest) ([]entity.Courier, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	filter := map[string]interface{}{
		"courier_type": input.CourierType,
		"status":       input.Status,
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

// swagger:route PUT /courier/{id} courier-update UpdateCourierRequest
// Update courier
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceImpl) UpdateCourier(uid string, input request.UpdateCourierRequest) message.Message {
	logger := log.With(s.logger, "CourierService", "UpdateCourier")

	_, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	data := map[string]interface{}{
		"name":         input.CourierName,
		"status":       input.Status,
		"logo":         input.Logo,
		"courier_type": input.CourierType,
	}

	err = s.courierRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.FailedMsg
}

// swagger:route DELETE /courier/courier/{id} courier-delete byParamDelete
// Delete courier
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceImpl) DeleteCourier(uid string) message.Message {
	logger := log.With(s.logger, "CourierService", "DeleteCourier")

	_, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	err = s.courierRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
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

	s.baseRepo.BeginTx()
	//Set request to entity
	CourierService := entity.CourierService{
		//General
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
		Status:              input.Status,

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
	}

	resultInsert, err := s.courierServiceRepo.CreateCourierService(&CourierService)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.ErrDB
	}
	s.baseRepo.CommitTx()

	return resultInsert, message.SuccessMsg
}

// swagger:route GET /courier/courier-services/{uid} Courier-Services CourierServiceRequestGetByUid
// Detail Courier Services
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) GetCourierService(uid string) (*entity.CourierService, message.Message) {
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

	return result, message.SuccessMsg
}

// swagger:route GET /courier/courier-services Courier-Services CourierServiceListRequest
// List of Courier Services
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceImpl) GetListCourierService(input request.CourierServiceListRequest) ([]entity.CourierService, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	filter := map[string]interface{}{
		"courier_uid":   input.CourierUID,
		"shipping_code": input.ShippingCode,
		"shipping_name": input.ShippingName,
		"status":        input.Status,
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
		return message.ErrNoData
	}

	err = s.courierServiceRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}
