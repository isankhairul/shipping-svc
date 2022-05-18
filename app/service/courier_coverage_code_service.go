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
	UpdateCourierCoverageCode(uid string, input request.SaveCourierCoverageCodeRequest) message.Message
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

// swagger:route POST /courier/courier-coverage-code  CreateCourierCoverageCode
// Create Courier Coverage Code
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse

func (s *CourierCoverageCodeServiceImpl) CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Create Courier Coverage Code")
	s.baseReo.BeginTx()
	var courier entity.Courier
	err := s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.ErrDB
	}

	// set request to entity
	courierCoverageCode := entity.CourierCoverageCode{
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
		return nil, message.ErrDB
	}
	s.baseReo.CommitTx()
	return result, message.SuccessMsg

}

// swagger:route GET /courier/courier-coverage-code/ CourierCoverageCodeList
// List products
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
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

// swagger:route GET /courier/courier-coverage-code/{id} GetCourierCoverageCode
// Get Courier
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *CourierCoverageCodeServiceImpl) GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Get Courier Coverage Code")

	result, err := s.courierCoverageCodeRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return result, message.SuccessMsg
}

// swagger:route PUT /courier/courier-coverage-code/{id}  UpdateCourierCoverageCode
// Update courier
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *CourierCoverageCodeServiceImpl) UpdateCourierCoverageCode(uid string, input request.SaveCourierCoverageCodeRequest) message.Message {
	logger := log.With(s.logger, "CourierCoverageCodeService", "List Courier Coverage Code")

	_, err := s.courierCoverageCodeRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	var courier entity.Courier
	err = s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
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

	err = s.courierCoverageCodeRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}
