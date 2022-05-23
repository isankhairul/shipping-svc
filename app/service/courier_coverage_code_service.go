package service

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"gorm.io/gorm"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type CourierCoverageCodeService interface {
	CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
	GetList(input request.CourierCoverageCodeListRequest) ([]entity.CourierCoverageCode, *base.Pagination, message.Message)
	GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message)
	UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
	ImportCourierCoverageCode(input request.ImportCourierCoverageCodeRequest) ([]response.ImportStatus, message.Message)
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
	count, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode, 0)

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
	result, err := s.courierCoverageCodeRepo.Create(&courierCoverageCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.FailedMsg
	}
	s.baseReo.CommitTx()
	courierCoverageCode.CourierUID = input.CourierUID
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

	for i := range result {
		var courier entity.Courier
		err := s.courierCoverageCodeRepo.GetCourierId(&courier, result[i].CourierID)
		if err != nil {
			return nil, nil, message.FailedMsg
		}
		result[i].CourierUID = courier.UID
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
func (s *CourierCoverageCodeServiceImpl) UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	s.baseReo.BeginTx()
	logger := log.With(s.logger, "CourierCoverageCodeService", "List Courier Coverage Code")

	courierCoverageCodeRepo, err := s.courierCoverageCodeRepo.FindByUid(input.Uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		return nil, message.ErrNoData
	}
	courierId := courierCoverageCodeRepo.ID
	var courier entity.Courier
	err = s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseReo.RollbackTx()
		mess := message.ErrDataExists
		mess.Message = "Not found Courier UID in Courier table"
		return nil, mess
	}

	var courierCoverageCode entity.CourierCoverageCode
	count, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode, courierId)

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
	result, err := s.courierCoverageCodeRepo.Update(input.Uid, data)
	if err != nil {
		s.baseReo.RollbackTx()
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}

	return result, message.SuccessMsg
}

// swagger:route POST /courier/courier-coverage-code/import/ Courier-Coverage-Code ImportCourierCoverageCodeRequest
// Import courier coverage code by CSV file
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *CourierCoverageCodeServiceImpl) ImportCourierCoverageCode(input request.ImportCourierCoverageCodeRequest) ([]response.ImportStatus, message.Message) {
	s.baseReo.BeginTx()
	logger := log.With(s.logger, "CourierCoverageCodeService", "Import Courier Coverage Codes")

	var resp []response.ImportStatus
	for _, row := range input.Rows {
		courierUid, courierUidOk := row["courier_uid"]
		countryCode, countryCodeOk := row["country_code"]
		postalCode, postalCodeOk := row["postal_code"]
		description, descriptionOk := row["description"]
		code1, code1Ok := row["code1"]
		code2, code2Ok := row["code2"]
		code3, code3Ok := row["code3"]
		code4, code4Ok := row["code4"]
		code5, code5Ok := row["code5"]
		code6, code6Ok := row["code6"]

		if !courierUidOk || !countryCodeOk || !postalCodeOk || !descriptionOk || !code1Ok || !code2Ok || !code3Ok || !code4Ok || !code5Ok || !code6Ok {
			return nil, message.ErrImportData
		}

		// Check empty string
		if courierUid == "" || countryCode == "" || postalCode == "" {
			resp = append(resp, response.ImportStatus{
				CourierUID:  courierUid,
				CountryCode: countryCode,
				PostalCode:  postalCode,
				Description: description,
				Code1:       code1,
				Code2:       code2,
				Code3:       code3,
				Code4:       code4,
				Code5:       code5,
				Code6:       code6,
				Message:     "Can not import missing Courier UID, Country Code, and Postal Code",
				Status:      false,
			})
			continue
		}

		// check courier_id existing
		var courier entity.Courier
		err := s.courierCoverageCodeRepo.GetCourierUid(&courier, courierUid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resp = append(resp, response.ImportStatus{
					CourierUID:  courierUid,
					CountryCode: countryCode,
					PostalCode:  postalCode,
					Description: description,
					Code1:       code1,
					Code2:       code2,
					Code3:       code3,
					Code4:       code4,
					Code5:       code5,
					Code6:       code6,
					Message:     "Not found Courier UID in Courier table",
					Status:      false,
				})
				continue
			} else {
				_ = level.Error(logger).Log(err)
				s.baseReo.RollbackTx()
				return nil, message.ErrDB
			}
		}
		// Check CourierUID, country_code and postal_code are unique
		var courierCoverageCode entity.CourierCoverageCode
		count, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, countryCode, postalCode, 0)

		if err != nil {
			_ = level.Error(logger).Log(err)
			s.baseReo.RollbackTx()
			return nil, message.ErrDB
		}

		if count == 0 {
			data := entity.CourierCoverageCode{
				CourierID:   courier.ID,
				CountryCode: countryCode,
				PostalCode:  postalCode,
				Description: description,
				Code1:       code1,
				Code2:       code2,
				Code3:       code3,
				Code4:       code4,
				Code5:       code5,
				Code6:       code6,
			}
			result, err := s.courierCoverageCodeRepo.Create(&data)
			if err != nil {
				_ = level.Error(logger).Log(err)
				s.baseReo.RollbackTx()
				return nil, message.ErrDB
			}
			resp = append(resp, response.ImportStatus{
				UID:         result.UID,
				CourierUID:  courierUid,
				CountryCode: countryCode,
				PostalCode:  postalCode,
				Description: description,
				Code1:       code1,
				Code2:       code2,
				Code3:       code3,
				Code4:       code4,
				Code5:       code5,
				Code6:       code6,
				Message:     "Created",
				Status:      true,
			})

		} else {
			data := map[string]interface{}{
				"courier_id":   courier.ID,
				"country_code": countryCode,
				"postal_code":  postalCode,
				"description":  description,
				"code1":        code1,
				"code2":        code2,
				"code3":        code3,
				"code4":        code4,
				"code5":        code5,
				"code6":        code6,
			}
			_, err := s.courierCoverageCodeRepo.Update(courierCoverageCode.UID, data)
			if err != nil {
				_ = level.Error(logger).Log(err)
				s.baseReo.RollbackTx()
				return nil, message.ErrDB
			}
			resp = append(resp, response.ImportStatus{
				UID:         courierCoverageCode.UID,
				CourierUID:  courierUid,
				CountryCode: countryCode,
				PostalCode:  postalCode,
				Description: description,
				Code1:       code1,
				Code2:       code2,
				Code3:       code3,
				Code4:       code4,
				Code5:       code5,
				Code6:       code6,
				Message:     "Updated",
				Status:      true,
			})
		}
	}
	s.baseReo.CommitTx()
	return resp, message.SuccessMsg
}
