package test

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var courierRepo = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
var courierService = service.NewCourierService(logger, baseRepository, courierRepo, courierServiceRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

func TestCreateCourier(t *testing.T) {
	req := request.SaveCourierRequest{
		CourierName: "test name",
		Code:        "test code",
		CourierType: "1",
		Logo:        "logo test",
	}
	courier := entity.Courier{}

	courierRepo.Mock.On("CreateCourier", &req).Return(&courier)
	courierRepo.Mock.On("FindByCode", mock.Anything).Return(nil)
	result, _ := courierService.CreateCourier(req)

	assert.NotNil(t, result)
	assert.Equal(t, "test name", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "test code", result.Code, "Code must be test code")
	assert.Equal(t, 1, result.Status, "Status must be 1")
	assert.Equal(t, "logo test", result.Logo, "Log  must be logo test")
}

func TestGetCourier(t *testing.T) {
	courier := entity.Courier{
		Code: "code test",
	}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	courierRepo.Mock.On("FindByUid", &uid).Return(courier)
	result, _ := courierService.GetCourier(uid)

	assert.NotNil(t, result, "Cannot nil")
	assert.Equal(t, "code test", result.Code, "Code must be code test")
}

func TestDeleteCourier(t *testing.T) {
	courier := entity.Courier{}

	uid := "BnOI8D7p9rR7tI1R9rySw"
	courierRepo.Mock.On("Delete", &uid).Return(courier)

	msg := courierService.DeleteCourier(uid)
	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

func TestListCouriers(t *testing.T) {
	req := request.CourierListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	couriers := []entity.Courier{
		{
			Code:        "test code",
			CourierName: "test name",
		},
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	courierRepo.Mock.On("FindByParams", 10, 1, "", mock.Anything).Return(couriers, &paginationResult)
	couriers, pagination, msg := courierService.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, 1, len(couriers), "Count of couriers must be 1")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}
