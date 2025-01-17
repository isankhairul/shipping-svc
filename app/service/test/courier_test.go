package test

import (
	//"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func init() {
// }
var courierstatus int32 = 1
var courierTest = entity.Courier{
	CourierName: "test name",
	Code:        "test code",
	CourierType: "1",
	Logo:        "logo test",
	Status:      &courierstatus,
}

func TestCreateCourier(t *testing.T) {
	req := request.SaveCourierRequest{
		CourierName: courierTest.CourierName,
		Code:        courierTest.Code,
		CourierType: courierTest.CourierType,
		Logo:        courierTest.Logo,
		Status:      *courierTest.Status,
	}
	courier := entity.Courier{}
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("CreateCourier", &req).Return(&courier)
	courierRepo.Mock.On("FindByCode", mock.Anything).Return(nil)
	var courierService = service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	result, _ := courierService.CreateCourier(req)

	assert.NotNil(t, result)
	assert.Equal(t, courierTest.CourierName, result.CourierName, courierNameIsNotCorrect)
	assert.Equal(t, courierTest.Code, result.Code, courierNameIsNotCorrect)
}

func TestGetCourier(t *testing.T) {
	courier := entity.Courier{
		Code: courierTest.Code,
	}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	var courierService = service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	result, _ := courierService.GetCourier(uid)

	assert.NotNil(t, result)
	assert.Equal(t, courierTest.Code, result.Code, codeIsNotCorrect)
}

func TestDeleteCourier(t *testing.T) {
	courier := entity.Courier{
		Code: courierTest.Code,
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{})
	courierRepo.Mock.On("Delete", mock.Anything).Return(nil)
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, codeIsNotCorrect)
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

func TestListCouriers(t *testing.T) {
	req := request.CourierListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	couriers := []response.CourierListResponse{
		{
			Code:        courierTest.Code,
			CourierName: courierTest.CourierName,
		},
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByParams", 10, 1, "", mock.Anything).Return(couriers, &paginationResult)
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	courierResponse, pagination, msg := courierService.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, codeIsNotCorrect)
	assert.Equal(t, 1, len(courierResponse), "Count of couriers must be 1")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}

func TestUpdateCourierFailNotFoundCourier(t *testing.T) {
	uid := "BnOI8D7p9rR7tI1R9rySw"
	err := errors.Wrap(errors.New("Courier found"), "courier")
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	courierRepo.Mock.On("FindByUid", mock.Anything).Return(nil, err)

	result, msg := courierService.UpdateCourier(uid, request.UpdateCourierRequest{})

	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrCourierNotFound.Code, codeIsNotCorrect)
}

func TestGetCourierFail(t *testing.T) {
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(nil, errors.New("Not found"))
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	result, msg := courierService.GetCourier(uid)

	assert.Nil(t, result)
	assert.Equal(t, msg.Code, message.ErrCourierNotFound.Code, codeIsNotCorrect)
}

func TestDeleteCourierFail(t *testing.T) {
	courier := entity.Courier{
		Code: courierTest.Code,
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{})
	courierRepo.Mock.On("Delete", mock.Anything).Return(errors.New("db"))
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrDB.Code, "message should be err db")
}

func TestDeleteCourierHasChildCourierService(t *testing.T) {
	courier := entity.Courier{
		BaseIDModel:     base.BaseIDModel{ID: 1},
		Code:            courierTest.Code,
		CourierServices: []*entity.CourierService{{}, {}},
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{CourierService: true})
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrCourierHasChildCourierService.Code, codeIsNotCorrect)
}

func TestDeleteCourierHasChildCourierCoverageCode(t *testing.T) {
	courier := entity.Courier{
		BaseIDModel:     base.BaseIDModel{ID: 1},
		Code:            courierTest.Code,
		CourierServices: []*entity.CourierService{{}, {}},
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{CourierCoverageCode: true})
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrCourierHasChildCourierCoverage.Code, codeIsNotCorrect)
}

func TestDeleteCourierHasChildChannelCourier(t *testing.T) {
	courier := entity.Courier{
		BaseIDModel:     base.BaseIDModel{ID: 1},
		Code:            courierTest.Code,
		CourierServices: []*entity.CourierService{{}, {}},
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{ChannelCourier: true})
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrCourierHasChildChannelCourier.Code, codeIsNotCorrect)
}

func TestDeleteCourierHasChildShippingCourierStatus(t *testing.T) {
	courier := entity.Courier{
		BaseIDModel:     base.BaseIDModel{ID: 1},
		Code:            courierTest.Code,
		CourierServices: []*entity.CourierService{{}, {}},
	}
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	courierRepo.Mock.On("IsCourierHasChild", courier.ID).Return(&entity.CourierHasChildFlag{ShippingCourierStatus: true})
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrCourierHasChildShippingStatus.Code, codeIsNotCorrect)
}

func TestDeleteCourierNotFound(t *testing.T) {
	uid := "BnOI8D7p9rR7tI1R9rySw"
	var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
	courierRepo.Mock.On("FindByUid", &uid).Return(nil)
	courierService := service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository, shipmentPredefinedRepository)
	msg := courierService.DeleteCourier(uid)

	assert.Equal(t, msg.Code, message.ErrCourierNotFound.Code, "message should be not found")
}
