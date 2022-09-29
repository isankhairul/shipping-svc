package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type CourierCoverageCodeService interface {
	CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
	GetList(input request.CourierCoverageCodeListRequest) ([]*entity.CourierCoverageCode, *base.Pagination, message.Message)
	GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message)
	UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message)
	DeleteCourierCoverageCode(uid string) message.Message
	ImportCourierCoverageCode(input request.ImportCourierCoverageCodeRequest) (*response.CourierCoverageCodeImportResponse, message.Message)
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

// swagger:operation POST /courier/courier-coverage-code/ Courier-Coverage-Code SaveCourierCoverageCodeRequest
// Create Courier Coverage Code
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/CourierCoverageCode'
func (s *CourierCoverageCodeServiceImpl) CreateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Create Courier Coverage Code")

	var courier entity.Courier
	err := s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrCourierNotFound
	}

	var courierCoverageCode entity.CourierCoverageCode
	count, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode, input.Subdistrict, 0)

	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrCourierCoverageCodeUidExist
	}

	if count > 0 {
		return nil, message.ErrCourierCoverageCodeExist
	}

	// set request to entity
	courierCoverageCode = entity.CourierCoverageCode{
		CourierID:   courier.ID,
		CountryCode: input.CountryCode,
		PostalCode:  input.PostalCode,
		Subdistrict: input.Subdistrict,
		Description: input.Description,
		Code1:       input.Code1,
		Code2:       input.Code2,
		Code3:       input.Code3,
		Code4:       input.Code4,
		Code5:       input.Code5,
		Code6:       input.Code6,
		Status:      &input.Status,
		BaseIDModel: base.BaseIDModel{
			CreatedBy: input.ActorName,
		},
	}
	result, err := s.courierCoverageCodeRepo.Create(&courierCoverageCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}
	courierCoverageCode.CourierUID = input.CourierUID
	return result, message.SuccessMsg

}

// swagger:operation GET /courier/courier-coverage-code/ Courier-Coverage-Code CourierCoverageCodeListRequest
// List couriers coverage code
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//           $ref: '#/definitions/MetaPaginationResponse'
//         data:
//           properties:
//             records:
//               type: array
//               items:
//                 $ref: '#/definitions/CourierCoverageCode'
func (s *CourierCoverageCodeServiceImpl) GetList(input request.CourierCoverageCodeListRequest) ([]*entity.CourierCoverageCode, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "List Courier Coverage Codes")

	filter := map[string]interface{}{
		"courier_name": input.Filters.CourierName,
		"country_code": input.Filters.CountryCode,
		"postal_code":  input.Filters.PostalCode,
		"subdistrict":  input.Filters.Subdistrict,
		"description":  input.Filters.Description,
		"status":       input.Filters.Status,
		"code1":        input.Filters.Code1,
		"code2":        input.Filters.Code2,
		"code3":        input.Filters.Code3,
		"code4":        input.Filters.Code4,
		"code5":        input.Filters.Code5,
		"code6":        input.Filters.Code6,
	}

	convertedFilterLower := make(map[string]interface{}, len(filter))
	for k, v := range filter {
		convertedFilterLower[strings.ToLower(k)] = v
	}

	result, pagination, err := s.courierCoverageCodeRepo.FindByParams(input.Limit, input.Page, input.Sort, convertedFilterLower)
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

// swagger:operation DELETE /courier/courier-coverage-code/{uid} Courier-Coverage-Code DeleteCourierCoverageCodeByIDParam
// Delete courier coverage code by UID
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           type: object
func (s *CourierCoverageCodeServiceImpl) DeleteCourierCoverageCode(uid string) message.Message {
	err := s.courierCoverageCodeRepo.DeleteByUid(uid)
	if err != nil {
		return message.ErrCourierCoverageCodeUidNotExist
	}
	return message.SuccessMsg
}

// swagger:operation GET /courier/courier-coverage-code/{uid} Courier-Coverage-Code CourierCoverageCodeByIDParam
// Get Courier Coverage Code by uid
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/CourierCoverageCode'
func (s *CourierCoverageCodeServiceImpl) GetCourierCoverageCode(uid string) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "Get Courier Coverage Code")

	result, err := s.courierCoverageCodeRepo.FindByUid(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrNoData
	}
	if result.Courier != nil {
		result.CourierUID = result.Courier.UID
	}
	return result, message.SuccessMsg
}

// swagger:operation PUT /courier/courier-coverage-code/{uid} Courier-Coverage-Code UpdateCourierCoverageCodeBody
// Update courier coverage by uid
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/CourierCoverageCode'
func (s *CourierCoverageCodeServiceImpl) UpdateCourierCoverageCode(input request.SaveCourierCoverageCodeRequest) (*entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "UpdateCourierCoverageCode")

	courierCoverageCodeRepo, err := s.courierCoverageCodeRepo.FindByUid(input.Uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrNoData
	}
	courierId := courierCoverageCodeRepo.ID

	var courier entity.Courier
	err = s.courierCoverageCodeRepo.GetCourierUid(&courier, input.CourierUID)
	if err != nil {
		_ = level.Error(logger).Log("s.courierCoverageCodeRepo.GetCourierUid", err.Error())
		return nil, message.ErrCourierNotFound
	}

	var courierCoverageCode entity.CourierCoverageCode
	count, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, input.CountryCode, input.PostalCode, input.Subdistrict, courierId)

	if err != nil {
		_ = level.Error(logger).Log("s.courierCoverageCodeRepo.CombinationUnique", err.Error())
		return nil, message.FailedMsg
	}

	if count > 0 {
		mess := message.ErrCourierCoverageCodeExist
		return nil, mess
	}

	data := map[string]interface{}{
		"courier_id":   courier.ID,
		"country_code": input.CountryCode,
		"postal_code":  input.PostalCode,
		"subdistrict":  input.Subdistrict,
		"description":  input.Description,
		"code1":        input.Code1,
		"code2":        input.Code2,
		"code3":        input.Code3,
		"code4":        input.Code4,
		"code5":        input.Code5,
		"code6":        input.Code6,
		"status":       input.Status,
		"updated_by":   input.ActorName,
	}
	result, err := s.courierCoverageCodeRepo.Update(input.Uid, data)
	if err != nil {
		_ = level.Error(logger).Log(s.courierCoverageCodeRepo.Update, err.Error())
		return nil, message.FailedMsg
	}
	result.CourierUID = input.CourierUID
	result.UID = input.Uid
	return result, message.SuccessMsg
}

// swagger:operation POST /courier/courier-coverage-code/import Courier-Coverage-Code ImportCourierCoverageCodeRequest
// Import courier coverage code by CSV file
//
// Description :
//
// ---
// consumes:
// - multipart/form-data
//
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/ImportCourierCoverageCode'
func (s *CourierCoverageCodeServiceImpl) ImportCourierCoverageCode(input request.ImportCourierCoverageCodeRequest) (*response.CourierCoverageCodeImportResponse, message.Message) {

	logger := log.With(s.logger, "CourierCoverageCodeService", "Import Courier Coverage Codes")
	totalRows := len(input.Rows)
	failedRows := 0
	successRows := 0

	data := []response.ImportStatus{}

	ok := s.checkImportedDataColumnValidity(input.Rows)
	if !ok {
		return nil, message.ErrImportData
	}

	for _, row := range input.Rows {
		countryCode := row["country_code"]
		postalCode := row["postal_code"]
		subdistrict := row["subdistrict"]
		description := row["description"]
		code1 := row["code1"]
		code2 := row["code2"]
		code3 := row["code3"]
		code4 := row["code4"]
		code5 := row["code5"]
		code6 := row["code6"]

		courier, failedData, failedRowCount := s.checkImportedDataRow(row)

		if failedRowCount > 0 {
			data = append(data, *failedData)
			failedRows++
			continue
		}

		// Check CourierUID, country_code and postal_code are unique
		var courierCoverageCode entity.CourierCoverageCode
		_, err := s.courierCoverageCodeRepo.CombinationUnique(&courierCoverageCode, courier.ID, countryCode, postalCode, subdistrict, 0)

		if err != nil {
			_ = level.Error(logger).Log(err)
			return nil, message.ErrDB
		}

		data := entity.CourierCoverageCode{
			CourierID:   courier.ID,
			CountryCode: countryCode,
			PostalCode:  postalCode,
			Subdistrict: subdistrict,
			Description: description,
			Code1:       code1,
			Code2:       code2,
			Code3:       code3,
			Code4:       code4,
			Code5:       code5,
			Code6:       code6,
			BaseIDModel: base.BaseIDModel{
				CreatedBy: input.ActorName,
			},
		}

		successCount, msg := s.upsert(courierCoverageCode.UID, data)
		if msg != message.SuccessMsg {
			return nil, msg
		}
		successRows += successCount
	}

	return &response.CourierCoverageCodeImportResponse{
		FailedData: data,
		Summary: response.ImportSummary{
			TotalRow:   totalRows,
			SuccessRow: successRows,
			FailedRow:  failedRows,
		},
	}, message.SuccessMsg
}

// return successCount, message
func (s *CourierCoverageCodeServiceImpl) upsert(uid string, input entity.CourierCoverageCode) (int, message.Message) {
	logger := log.With(s.logger, "CourierCoverageCodeService", "upsert")
	var log string
	var err error

	if len(uid) == 0 {
		log = "s.courierCoverageCodeRepo.Create"
		_, err = s.courierCoverageCodeRepo.Create(&input)

	} else {
		log = "s.courierCoverageCodeRepo.Update"
		data := map[string]interface{}{
			"courier_id":   input.CourierID,
			"country_code": input.CountryCode,
			"postal_code":  input.PostalCode,
			"subdistrict":  input.Subdistrict,
			"description":  input.Description,
			"code1":        input.Code1,
			"code2":        input.Code2,
			"code3":        input.Code3,
			"code4":        input.Code4,
			"code5":        input.Code5,
			"code6":        input.Code6,
			"updated_by":   input.CreatedBy,
		}
		_, err = s.courierCoverageCodeRepo.Update(uid, data)

	}

	if err != nil {
		_ = level.Error(logger).Log(log, err.Error())
		return 0, message.ErrDB
	}

	return 1, message.SuccessMsg
}

func (s *CourierCoverageCodeServiceImpl) checkImportedDataColumnValidity(input []map[string]string) bool {
	var row map[string]string

	if len(input) == 0 {
		return true
	}
	row = input[0]

	_, courierUidOk := row["courier_uid"]
	_, countryCodeOk := row["country_code"]
	_, postalCodeOk := row["postal_code"]
	_, subdistrictOk := row["subdistrict"]
	_, descriptionOk := row["description"]
	_, code1Ok := row["code1"]
	_, code2Ok := row["code2"]
	_, code3Ok := row["code3"]
	_, code4Ok := row["code4"]
	_, code5Ok := row["code5"]
	_, code6Ok := row["code6"]

	return courierUidOk && countryCodeOk && postalCodeOk && subdistrictOk && descriptionOk && code1Ok && code2Ok && code3Ok && code4Ok && code5Ok && code6Ok
}

// return courier, array failed data, failedRowCount, summaryRowCount, message
func (s *CourierCoverageCodeServiceImpl) checkImportedDataRow(row map[string]string) (*entity.Courier, *response.ImportStatus, int) {
	courierUid := row["courier_uid"]
	countryCode := row["country_code"]
	postalCode := row["postal_code"]
	subdistrict := row["subdistrict"]
	description := row["description"]
	code1 := row["code1"]
	code2 := row["code2"]
	code3 := row["code3"]
	code4 := row["code4"]
	code5 := row["code5"]
	code6 := row["code6"]

	msg := "Can not import missing Courier UID, Country Code, and Postal Code"

	if courierUid != "" && countryCode != "" && postalCode != "" {

		// check courier_uid is exist
		var courier entity.Courier
		err := s.courierCoverageCodeRepo.GetCourierUid(&courier, courierUid)

		if err == nil {
			return &courier, nil, 0
		}

		msg = message.ErrCourierNotFound.Message
	}

	// check courier_id existing
	return nil, &response.ImportStatus{
		CourierUID:  courierUid,
		CountryCode: countryCode,
		PostalCode:  postalCode,
		Subdistrict: subdistrict,
		Description: description,
		Code1:       code1,
		Code2:       code2,
		Code3:       code3,
		Code4:       code4,
		Code5:       code5,
		Code6:       code6,
		Status:      false,
		Message:     msg,
	}, 1
}
