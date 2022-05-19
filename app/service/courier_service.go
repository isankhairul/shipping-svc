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

type CourierService interface {
	CreateCourier(input request.SaveCourierRequest) (*entity.Courier, message.Message)
	GetList(input request.CourierListRequest) ([]entity.Courier, *base.Pagination, message.Message)
	UpdateCourier(uid string, input request.SaveCourierRequest) message.Message
	GetCourier(uid string) (*entity.Courier, message.Message)
	DeleteCourier(uid string) message.Message
}

type courierServiceImpl struct {
	logger      log.Logger
	baseRepo    repository.BaseRepository
	courierRepo repository.CourierRepository
}

func NewCourierService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.CourierRepository,
) CourierService {
	return &courierServiceImpl{lg, br, pr}
}

// swagger:route POST /courier/courier Courier ManagedCourierRequest
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
func (s *courierServiceImpl) UpdateCourier(uid string, input request.SaveCourierRequest) message.Message {
	logger := log.With(s.logger, "CourierService", "UpdateCourier")

	_, err := s.courierRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	data := map[string]interface{}{
		"courier_name": input.CourierName,
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
