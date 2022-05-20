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

type CourierCoverageCodeService interface {
	CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
	GetList(input request.CourierCoverageCodeListRequest) ([]entity.CourierCoverageCode, *base.Pagination, message.Message)
	GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message)
	UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) message.Message
}

type CourierCoverageCodeServiceImpl struct {
	logger                  log.Logger
	baseReo                 repository.BaseRepository
	courierCoverageCodeRepo repository.CourierCoverageCodeRepository
}

func NewCourierCoverageCodeService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.CourierCoverageCodeRepository,
) CourierCoverageCodeService {
	return &CourierCoverageCodeServiceImpl{lg, br, pr}
}

// swagger:route POST /courier/courier-coverage-code/ Courier-Coverage-Code SaveCourierCoverageCodeRequest
// Create Courier Coverage Code
//
// responses:
//  401: errorResponse
//  201: CourierCoverageCode

func (s *CourierCoverageCodeServiceImpl) CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Create Courier Coverage Code")
	s.baseReo.BeginTx()

	var courier entity.Courier
	err := s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		mess := message.ErrDataExists
		mess.Message = "Not found courier_uid"
		return nil, mess
	}
	var courierCoverageCode entity.CourierCoverageCode
	count, err := s.courierCoverageCodeRepo.CombinationUnique(courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode)

	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.FailedMsg
	}

	if count > 0 {
		mess := message.Message{Code: 34005, Message: "The combination of courier_uid, country_code and postal_code is exist in database"}
		s.baseReo.RollbackTx()
		return nil, mess
	}

	// set request to entity
	courierCoverageCode = entity.CourierCoverageCode{
		CourierID:   courier.ID,
		CountryCode: input.CountryCode,
		PostalCode:  input.PostalCode,
		Description: input.Description,
		Code1:       input.Code1,
		Code2:       input.Code2,
		Code3:       input.Code3,
		Code4:       input.Code3,
		Code5:       input.Code5,
		Code6:       input.Code6,
	}
	result, err := s.courierCoverageCodeRepo.CreateCourierCoverageCodeRepo(&courierCoverageCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.FailedMsg
	}
	s.baseReo.CommitTx()
	result.CourierUID = courier.UID
	return result, message.SuccessMsg

}

// swagger:route GET /courier/courier-coverage-code/ Courier-Coverage-Code CourierCoverageCodeListRequest
// List products
//
// responses:
//  401: SuccessResponse
//  200: PaginationResponse
func (s *CourierCoverageCodeServiceImpl) GetList(input request.CourierCoverageCodeListRequest) ([]entity.CourierCoverageCode, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "List Courier Coverage Codes")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	result, pagination, err := s.courierCoverageCodeRepo.FindByParams(input.Limit, input.Page, input.Sort)
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

// swagger:route GET /courier/courier-coverage-code/{uid} Courier-Coverage-Code CourierCoverageCodeByIDParam
// Get Courier Coverage Code by uid
//
// responses:
//  401: SuccessResponse
//  200: CourierCoverageCode
func (s *CourierCoverageCodeServiceImpl) GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Get Courier Coverage Code")

	result, err := s.courierCoverageCodeRepo.FindByUid(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrNoData
	}

	var courier entity.Courier
	err = s.courierCoverageCodeRepo.GetCourierId(&courier, result.CourierID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.ErrDB
	}
	result.CourierUID = courier.UID
	return result, message.SuccessMsg
}

// swagger:route PUT /courier/courier-coverage-code/{uid} Courier-Coverage-Code SaveCourierCoverageCodeRequest
// Update courier coverage by uid
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *CourierCoverageCodeServiceImpl) UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) message.Message {
	s.baseReo.BeginTx()
	logger := log.With(s.logger, "CourierCoverageCodeService", "List Courier Coverage Code")

	_, err := s.courierCoverageCodeRepo.FindByUid(input.Uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return message.ErrNoData
	}

	var courier entity.Courier
	err = s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		mess := message.ErrDataExists
		mess.Message = "Not found courier_uid"
		return mess
	}

	var courierCoverageCode entity.CourierCoverageCode
	count, err := s.courierCoverageCodeRepo.CombinationUnique(courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode)

	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return message.FailedMsg
	}

	if count > 0 {
		mess := message.Message{Code: 34005, Message: "The combination of courier_uid, country_code and postal_code is exist in database"}
		s.baseReo.RollbackTx()
		return mess
	}

	data := map[string]interface{}{
		"courier_id":   courier.ID,
		"country_code": input.CountryCode,
		"postal_code":  input.PostalCode,
		"description":  input.Description,
		"code1":        input.Code1,
		"code2":        input.Code2,
		"code3":        input.Code3,
		"code4":        input.Code3,
		"code5":        input.Code5,
		"code6":        input.Code6,
	}

	err = s.courierCoverageCodeRepo.Update(input.Uid, data)
	if err != nil {
		s.baseReo.RollbackTx()
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}
