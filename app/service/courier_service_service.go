package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type CourierServiceService interface {
	CreateCourierService(input request.SaveCourierServiceRequest) (*entity.CourierService, message.Message)
	GetList(input request.CourierServiceListRequest) ([]entity.CourierService, *base.Pagination, message.Message)
	UpdateCourierService(uid string, input request.SaveCourierServiceRequest) message.Message
	GetCourierService(uid string) (*entity.CourierService, message.Message)
	DeleteCourierService(uid string) message.Message
}

type courierServiceServiceImpl struct {
	logger             log.Logger
	baseRepo           repository.BaseRepository
	courierServiceRepo repository.CourierServiceRepository
}

func NewCourierServiceService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.CourierServiceRepository,
) CourierServiceService {
	return &courierServiceServiceImpl{lg, br, pr}
}

// swagger:route POST /courier/courier-services Courier Service ManagedCourierServiceRequest
// Manage CourierService
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceServiceImpl) CreateCourierService(input request.SaveCourierServiceRequest) (*entity.CourierService, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "CreateCourierService")
	//Check exits `courier_id & shiping_code`
	//Set default value
	defaultLimit := 10
	defaultPage := 1
	defaultSort := ""
	filter := map[string]interface{}{
		"shipping_code": input.ShippingCode,
		"courier_id":    input.CourierId,
	}

	result, _, err := s.courierServiceRepo.FindByParams(defaultLimit, defaultPage, defaultSort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}

	if result != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDataCourierServiceExists
	}

	s.baseRepo.BeginTx()
	//Set request to entity
	CourierService := entity.CourierService{
		//General
		CourierId:           input.CourierId,
		CourierName:         input.CourierName,
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

// swagger:route GET /courier/courier-services/{uid} Get-Courier-Service CourierService
// Get CourierService
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *courierServiceServiceImpl) GetCourierService(uid string) (*entity.CourierService, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "GetCourierService")

	result, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoDataCourierService
	}

	return result, message.SuccessMsg
}

func (s *courierServiceServiceImpl) GetList(input request.CourierServiceListRequest) ([]entity.CourierService, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierServiceService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	filter := map[string]interface{}{
		"shipping_type": input.ShippingType,
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

// swagger:route PUT courier/courier-services/{uid} courier-service-update UpdateCourierServiceRequest
// Update courierservice
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceServiceImpl) UpdateCourierService(uid string, input request.SaveCourierServiceRequest) message.Message {
	logger := log.With(s.logger, "CourierServiceService", "UpdateCourierService")

	_, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	data := map[string]interface{}{
		"courier_id":           input.CourierId,
		"courier_name":         input.CourierName,
		"shipping_code":        input.ShippingCode,
		"shipping_name":        input.ShippingName,
		"shipping_type":        input.ShippingType,
		"shipping_description": input.ShippingDescription,
		"ETD_Min":              input.ETD_Min,
		"ETD_Max":              input.ETD_Max,
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
		return message.FailedMsg
	}

	return message.FailedMsg
}

// swagger:route DELETE /courier/courier-services/{uid} courier-delete byParamDelete
// Delete courierservice
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *courierServiceServiceImpl) DeleteCourierService(uid string) message.Message {
	logger := log.With(s.logger, "CourierServiceService", "DeleteCourierService")

	_, err := s.courierServiceRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	err = s.courierServiceRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}
