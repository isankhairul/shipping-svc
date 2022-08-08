package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ShipmentPredefinedService interface {
	GetAll(input request.ListShipmentPredefinedRequest) ([]*entity.ShippmentPredefined, *base.Pagination, message.Message)
	UpdateShipmentPredefined(input request.UpdateShipmentPredefinedRequest) (*entity.ShippmentPredefined, message.Message)
	CreateShipmentPredefined(input request.CreateShipmentPredefinedRequest) (*entity.ShippmentPredefined, message.Message)
	GetByUID(uid string) (*response.ShippmentPredefined, message.Message)
}

type ShipmentPredefinedServiceImpl struct {
	logger     log.Logger
	baseRepo   repository.BaseRepository
	predefines repository.ShipmentPredefinedRepository
}

func NewShipmentPredefinedService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.ShipmentPredefinedRepository,
) ShipmentPredefinedService {
	return &ShipmentPredefinedServiceImpl{lg, br, pr}
}

func (s *ShipmentPredefinedServiceImpl) CreateShipmentPredefined(input request.CreateShipmentPredefinedRequest) (*entity.ShippmentPredefined, message.Message) {
	logger := log.With(s.logger, "ShipmentPredefinedService", "Create")
	var count int64
	_ = s.baseRepo.GetDB().Model(&entity.ShippmentPredefined{}).Count(&count)
	if count > 0 {
		return nil, message.ErrDataExists
	}
	predefined := &entity.ShippmentPredefined{Type: input.Type, Title: input.Type, Code: input.Code, Note: input.Note, Status: 1}
	err := s.baseRepo.GetDB().Create(&predefined).Error
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}

	return predefined, message.SuccessMsg
}

// swagger:route PUT /other/shipment-predefined/{uid} Courier-Predefined UpdateShipmentPredefinedRequest
// Update a shipment predefined
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *ShipmentPredefinedServiceImpl) UpdateShipmentPredefined(input request.UpdateShipmentPredefinedRequest) (*entity.ShippmentPredefined, message.Message) {
	logger := log.With(s.logger, "ShipmentPredefinedService", "Update")

	var ret *entity.ShippmentPredefined

	_, err := s.predefines.GetShipmentPredefinedByUid(input.Uid)

	if err != nil {
		return nil, message.ErrShipmentPredefinedNotFound
	}

	predefined := &entity.ShippmentPredefined{Type: input.Type,
		Title: input.Title,
		Code:  input.Code, Status: input.Status,
		Note: input.Note, BaseIDModel: base.BaseIDModel{UID: input.Uid}}
	ret, err = s.predefines.UpdateShipmentPredefined(*predefined)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}

	return ret, message.SuccessMsg
}

// swagger:route GET /other/shipment-predefined Courier-Predefined ListShipmentPredefinedRequest
// Get predefined
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ShipmentPredefinedServiceImpl) GetAll(input request.ListShipmentPredefinedRequest) ([]*entity.ShippmentPredefined, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ShipmentPredefinedService", "GetAll")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	filters := map[string]interface{}{
		"code":   input.Filters.Code,
		"title":  input.Filters.Title,
		"type":   input.Filters.Type,
		"status": input.Filters.Status,
	}
	items, pagination, err := s.predefines.GetAll(input.Limit, input.Page, input.Sort, filters)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, pagination, message.ErrNoData
	}
	return items, pagination, message.SuccessMsg
}

// swagger:route GET /other/shipment-predefined/{uid} Courier-Predefined GetShipmentPredefinedByUID
// Get shipment predefined by UID
//
// responses:
//  200: GetShippmentPredefined
func (s *ShipmentPredefinedServiceImpl) GetByUID(uid string) (*response.ShippmentPredefined, message.Message) {
	logger := log.With(s.logger, "ShipmentPredefinedService", "GetByUID")

	var result *entity.ShippmentPredefined

	result, err := s.predefines.GetShipmentPredefinedByUid(uid)

	if err != nil {
		_ = level.Error(logger).Log("error", err.Error())
		return nil, message.ErrShipmentPredefinedNotFound
	}

	if result == nil {
		return nil, message.ErrShipmentPredefinedNotFound
	}

	return &response.ShippmentPredefined{
		UID:    result.UID,
		Type:   result.Type,
		Code:   result.Code,
		Title:  result.Title,
		Note:   result.Note,
		Status: result.Status,
	}, message.SuccessMsg
}
