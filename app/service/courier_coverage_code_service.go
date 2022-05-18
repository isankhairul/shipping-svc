package service

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type CourierCoverageCodeService interface {
	CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
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

// swagger:route POST /courier/courier-coverage-code CourierCoverageCode SaveCourierCoverageCodeRequest
// Create Courier Coverage Code
//
// responses:
// 200:
// 	description: Success
// 	headers:
// 	X-Correlation-ID:
// 		$ref: '#/components/headers/XCorrelationID'
// 	content:
// 	application/json:
// 		schema:
// 		type: object
// 		properties:
// 			meta:
// 			$ref: '#/components/schemas/Meta200'
// 			data:
// 			type : object
// 			properties:
// 				records:
// 				$ref: '#/components/schemas/CourierCoverageCode'
// 			errors:
// 			type: object
// 			example: {}
// 401:
// 	$ref: '#/components/responses/Unauthorized'
// 400:
// 	$ref: '#/components/responses/InvalidRequestData'
// 500:
// 	$ref: '#/components/responses/InternalServerError'

func (s *CourierCoverageCodeServiceImpl) CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Create Courier Coverage Code")
	s.baseReo.BeginTx()
	// set request to entity
	courierCoverageCode := entity.CourierCoverageCode{
		// CourierId: input.CourierUID,
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
	result, err := s.courierCoverageCodeRepo.CreateCoverageCodeRepo(&courierCoverageCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.ErrDB
	}
	s.baseReo.CommitTx()
	return result, message.SuccessMsg

}
