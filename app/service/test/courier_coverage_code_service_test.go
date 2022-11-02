package test

import (
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

// func init() {
// }

var vn = entity.CourierCoverageCode{
	CourierUID:          "UCMvWngocMqKbaC3AWQBF",
	CountryCode:         "VN",
	ProvinceNumericCode: "1",
	ProvinceName:        "Province A",
	CityNumericCode:     "12",
	CityName:            "City A",
	PostalCode:          "70000",
	Description:         "Vietnam code",
}

func TestCreateCourierCoverageCode(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	req := request.SaveCourierCoverageCodeRequest{
		CourierUID:          vn.CourierUID,
		CountryCode:         vn.CountryCode,
		ProvinceNumericCode: vn.ProvinceNumericCode,
		ProvinceName:        vn.ProvinceName,
		CityNumericCode:     vn.CityNumericCode,
		CityName:            vn.CityName,
		PostalCode:          vn.PostalCode,
		Description:         vn.Description,
	}

	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(nil)
	courierCoverageCodeRepository.Mock.On("CombinationUnique", mock.Anything).Return(0, nil)

	result, _ := svcCourierCoverageCode.CreateCourierCoverageCode(req)
	assert.NotNil(t, result)
	assert.Equal(t, vn.CourierUID, result.CourierUID, courierUIDIsNotCorrect)
	assert.Equal(t, vn.CountryCode, result.CountryCode, uidIsNotCorrect)
	assert.Equal(t, vn.ProvinceNumericCode, result.ProvinceNumericCode, "Province numeric code is not correct")
	assert.Equal(t, vn.ProvinceName, result.ProvinceName, "Province name is not correct")
	assert.Equal(t, vn.CityNumericCode, result.CityNumericCode, "City numeric code is not correct")
	assert.Equal(t, vn.CityName, result.CityName, "City name is not correct")
}

func TestUpdateCourierCoverageCode(t *testing.T) {
	courierCoverageCode := vn

	req := request.SaveCourierCoverageCodeRequest{
		Uid:         "123",
		CourierUID:  vn.CourierUID,
		CountryCode: vn.CountryCode,
		PostalCode:  vn.PostalCode,
		Description: vn.Description,
	}
	var courier entity.Courier
	courier.UID = vn.UID
	courier.ID = 555

	courierCoverageCodeRepository.Mock.On("FindByUid", req.Uid).Return(courierCoverageCode)
	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(courier)
	courierCoverageCodeRepository.Mock.On("CombinationUnique", mock.Anything).Return(0, nil)
	courierCoverageCodeRepository.Mock.On("Update", req.Uid, mock.Anything).Return(&courierCoverageCode, nil)

	result, _ := svcCourierCoverageCode.UpdateCourierCoverageCode(req)
	assert.NotNil(t, result)
	assert.Equal(t, vn.CourierUID, result.CourierUID, courierUIDIsNotCorrect)
	assert.Equal(t, vn.CountryCode, result.CountryCode, uidIsNotCorrect)
}

func TestGetCourierCoverageCode(t *testing.T) {
	var courier entity.Courier
	courier.UID = vn.UID
	courier.ID = 555
	uid := "123"

	courierCoverageCodeRepository.Mock.On("FindByUid", uid).Return(vn)
	// courierCoverageCodeRepository.Mock.On("GetCourierId", mock.Anything, mock.Anything).Return(courier)

	result, _ := svcCourierCoverageCode.GetCourierCoverageCode(uid)
	assert.NotNil(t, result)
	assert.Equal(t, vn.CourierUID, result.CourierUID, courierUIDIsNotCorrect)
	assert.Equal(t, vn.CountryCode, result.CountryCode, uidIsNotCorrect)
}

func TestDeleteCourierCoverageCode(t *testing.T) {
	uid := vn.UID

	courierCoverageCodeRepository.Mock.On("DeleteByUid", mock.Anything).Return(nil)

	message := svcCourierCoverageCode.DeleteCourierCoverageCode(uid)
	assert.NotNil(t, message)
	assert.NotNil(t, 201000, message.Code, "Expected 201000")
}

func TestListCourierCoverageCode(t *testing.T) {
	req := request.CourierCoverageCodeListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	courierCoverageCode := []*entity.CourierCoverageCode{
		{
			CourierUID:  vn.CourierUID,
			CountryCode: vn.CountryCode,
			PostalCode:  vn.PostalCode,
			Description: vn.Description,
		},
		{
			CourierUID:  "UCMvWngocMqKbaC3AWQBF2",
			CountryCode: "US",
			PostalCode:  "10000",
			Description: "US code",
		},
		{
			CourierUID:  "UCMvWngocMqKbaC3AWQBF3",
			CountryCode: "UK",
			PostalCode:  "20000",
			Description: "UK code",
		},
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	courierCoverageCodeRepository.Mock.On("FindByParams", 10, 1, "").Return(courierCoverageCode, &paginationResult)
	courierCoverageCodeRepository.Mock.On("GetCourierId", mock.Anything, mock.Anything).Return(entity.Courier{})
	courierCoverageCodes, pagination, msg := svcCourierCoverageCode.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, 3, len(courierCoverageCodes), "Count of courier must be 3")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")
}

func TestImportCourierCoverageCode(t *testing.T) {
	req := request.ImportCourierCoverageCodeRequest{
		Rows: []map[string]string{
			{
				"courier_uid":           vn.CourierUID,
				"country_code":          vn.CountryCode,
				"province_numeric_code": vn.ProvinceNumericCode,
				"province_name":         vn.ProvinceName,
				"city_numeric_code":     vn.CityNumericCode,
				"city_name":             vn.CityName,
				"postal_code":           vn.PostalCode,
				"district_numeric_code": vn.DistrictNumericCode,
				"district_name":         vn.DistrictName,
				"subdistrict":           vn.Subdistrict,
				"subdistrict_name":      vn.SubdistrictName,
				"description":           vn.Description,
				"code1":                 "",
				"code2":                 "",
				"code3":                 "",
				"code4":                 "",
				"code5":                 "",
				"code6":                 "",
			},
			{
				"courier_uid":           "UCMvWngocMqKbaC3AWQBF",
				"country_code":          "",
				"province_numeric_code": "",
				"province_name":         "",
				"city_numeric_code":     "",
				"city_name":             "",
				"postal_code":           "",
				"subdistrict":           vn.Subdistrict,
				"description":           vn.Description,
				"code1":                 "",
				"code2":                 "",
				"code3":                 "",
				"code4":                 "",
				"code5":                 "",
				"code6":                 "",
			},
		},
	}
	courierCoverageCodeRepository.Mock.On("Update", mock.Anything, mock.Anything).Return(&entity.CourierCoverageCode{}, nil)
	result, _ := svcCourierCoverageCode.ImportCourierCoverageCode(req)
	data := result.FailedData
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(data), "Count of result must be 1")
	assert.Equal(t, vn.CourierUID, data[0].CourierUID, courierUIDIsNotCorrect)
}

func TestCreateCourierCoverageCodeFailedWithDuplicatedUniqueCode(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	req := request.SaveCourierCoverageCodeRequest{
		CourierUID:  vn.CourierUID,
		CountryCode: vn.CountryCode,
		PostalCode:  vn.PostalCode,
		Description: vn.Description,
	}

	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(nil)
	courierCoverageCodeRepository.Mock.On("CombinationUnique", mock.Anything).Return(0, errors.New("Found"))

	result, msg := svcCourierCoverageCode.CreateCourierCoverageCode(req)

	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrCourierCoverageCodeUidExist.Code, "Duplicated coverage code with courier")
}

func TestUpdateCourierCoverageCodefailedWithNotFound(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	req := request.SaveCourierCoverageCodeRequest{
		Uid:         "123",
		CourierUID:  vn.CourierUID,
		CountryCode: vn.CountryCode,
		PostalCode:  vn.PostalCode,
		Description: vn.Description,
	}
	var courier entity.Courier
	courier.UID = "UCMvWngocMqKbaC3AWQBF"
	courier.ID = 555

	courierCoverageCodeRepository.Mock.On("FindByUid", mock.Anything).Return(nil, errors.New("Not Found"))
	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(errors.New("Found"))
	result, msg := svcCourierCoverageCode.UpdateCourierCoverageCode(req)

	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrNoData.Code, "Not found courier coverage by uid")
}

func TestUpdateCourierCoverageCodefailedWithDuplicated(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	req := request.SaveCourierCoverageCodeRequest{
		Uid:         "123",
		CourierUID:  vn.CourierUID,
		CountryCode: vn.CountryCode,
		PostalCode:  vn.PostalCode,
		Description: vn.Description,
	}

	courierCoverageCode := entity.CourierCoverageCode{
		CourierID:   555,
		CourierUID:  vn.CourierUID,
		CountryCode: vn.CountryCode,
		PostalCode:  vn.PostalCode,
		Description: vn.Description,
		Courier: &entity.Courier{
			BaseIDModel: base.BaseIDModel{UID: "UCMvWngocMqKbaC3AWQBF"},
		},
	}
	var courier entity.Courier
	courier.UID = "UCMvWngocMqKbaC3AWQBF"
	courier.ID = 555

	courierCoverageCodeRepository.Mock.On("FindByUid", mock.Anything).Return(courierCoverageCode, nil)
	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(courier)
	courierCoverageCodeRepository.Mock.On("CombinationUnique", mock.Anything).Return(0, errors.New("Found"))
	result, msg := svcCourierCoverageCode.UpdateCourierCoverageCode(req)

	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.FailedMsg.Code, "Duplicated coverage code")
}

func TestGetCourierCoverageCodeFailed(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	courierCoverageCodeRepository.Mock.On("FindByUid", mock.Anything).Return(nil, errors.New("Not found"))
	result, msg := svcCourierCoverageCode.GetCourierCoverageCode("123")
	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrNoData.Code, "Coverage not exists")
}

func TestDeleteCourierCoverageCodeFailed(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	courierCoverageCodeRepository.Mock.On("DeleteByUid", mock.Anything).Return(errors.New("Not found"))
	msg := svcCourierCoverageCode.DeleteCourierCoverageCode("123")
	assert.Equal(t, msg.Code, message.ErrCourierCoverageCodeUidNotExist.Code, "Coverage not exists")
}

func TestImportCourierCoverageCodeFailed(t *testing.T) {
	req := request.ImportCourierCoverageCodeRequest{
		Rows: []map[string]string{
			{
				"courier_uid":  vn.CourierUID,
				"country_code": vn.CountryCode,
				"description":  vn.Description,
				"code1":        "",
				"code2":        "",
				"code3":        "",
				"code4":        "",
				"code5":        "",
				"code6":        "",
			},
		},
	}
	courierCoverageCodeRepository.Mock.On("Update", mock.Anything, mock.Anything).Return(&entity.CourierCoverageCode{}, nil)
	result, msg := svcCourierCoverageCode.ImportCourierCoverageCode(req)
	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrImportData.Code)
}

func TestImportCourierCoverageCodeFailedWithNotFoundCourier(t *testing.T) {
	var baseCourierCoverageCodeRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
	var courierCoverageCodeRepository = &repository_mock.CourierCoverageCodeRepositoryMock{Mock: mock.Mock{}}
	var svcCourierCoverageCode = service.NewCourierCoverageCodeService(logger, baseCourierCoverageCodeRepository, courierCoverageCodeRepository)

	req := request.ImportCourierCoverageCodeRequest{
		Rows: []map[string]string{
			{
				"courier_uid":           vn.CourierUID,
				"country_code":          vn.CountryCode,
				"province_numeric_code": vn.ProvinceNumericCode,
				"province_name":         vn.ProvinceName,
				"city_numeric_code":     vn.CityNumericCode,
				"city_name":             vn.CityName,
				"description":           vn.Description,
				"postal_code":           "any",
				"district_numeric_code": vn.DistrictNumericCode,
				"district_name":         vn.DistrictName,
				"subdistrict":           vn.Subdistrict,
				"subdistrict_name":      vn.SubdistrictName,
				"code1":                 "",
				"code2":                 "",
				"code3":                 "",
				"code4":                 "",
				"code5":                 "",
				"code6":                 "",
			},
		},
	}
	courierCoverageCodeRepository.Mock.On("Update", mock.Anything, mock.Anything).Return(&entity.CourierCoverageCode{}, nil)
	courierCoverageCodeRepository.Mock.On("GetCourierUid", mock.Anything).Return(errors.New("Found"))
	result, msg := svcCourierCoverageCode.ImportCourierCoverageCode(req)
	assert.NotNil(t, result)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code)
}
