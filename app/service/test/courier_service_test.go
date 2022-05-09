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

// var logger log.Logger

// var baseRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var courierRepository = &repository_mock.CourierRepositoryMock{Mock: mock.Mock{}}
var courierSvc = service.NewCourierService(logger, baseRepository, courierRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	// db.AutoMigrate(&entity.Courier{})
}

func TestCreateCourier(t *testing.T) {
	req := request.SaveCourierRequest{
		CourierName: "Prenagen",
		Code:        "123",
		CourierType: "SKU_X",
		Logo:        "Pcs",
	}

	result, _ := courierSvc.CreateCourier(req)

	assert.NotNil(t, result)
	assert.Equal(t, "Prenagen", result.CourierName, "Courier Name must be Prenagen")
	assert.Equal(t, "123", result.Code, "Courier code must be 123")
	assert.Equal(t, "SKU_X", result.CourierType, "Type must be SKU_X")
	assert.Equal(t, "Pcs", result.Logo, "UOM must be Pcs")
}

func TestGetCourier(t *testing.T) {
	courier := entity.Courier{
		CourierName: "Prenagen",
		Code:        "12345",
		CourierType: "CourierType",
		Logo:        "Pcs",
	}

	uid := "123"
	courierRepository.Mock.On("FindByUid", &uid).Return(courier)
	result, _ := courierSvc.GetCourier(uid)

	assert.NotNil(t, result, "Cannot nil")
	assert.Equal(t, "Prenagen", result.CourierName, "Name must be Prenagen")
	assert.Equal(t, "12345", result.Code, "Weight must be 12345")
	assert.Equal(t, "CourierType", result.CourierType, "Courier Type must be CourierType")
	assert.Equal(t, "Pcs", result.Logo, "Logo must be Pcs")
}

func TestDeleteCourier(t *testing.T) {
	courier := entity.Courier{
		CourierName: "Prenagen",
		Code:        "12345",
		CourierType: "SKU_X",
		Logo:        "Pcs",
	}

	uid := "123"
	courierRepository.Mock.On("FindByUid", &uid).Return(courier)
	msg := courierSvc.DeleteCourier(uid)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be Null")
}

func TestListCourier(t *testing.T) {
	req := request.CourierListRequest{
		Page:  1,
		Sort:  "",
		Limit: 10,
	}

	courier := []entity.Courier{
		{
			CourierName: "Prenagen",
			Code:        "123",
			CourierType: "SKU_X",
			Logo:        "Pcs",
		},
		{
			CourierName: "Prenage",
			Code:        "1234",
			CourierType: "SYU_X",
			Logo:        "Pcs",
		},
		{
			CourierName: "Prenag",
			Code:        "12345",
			CourierType: "SKUU_X",
			Logo:        "Pcs",
		},
	}

	filter := map[string]interface{}{
		"courier_type": "",
		"status":       0,
	}

	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	courierRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(courier, &paginationResult)

	couriers, pagination, msg := courierSvc.GetList(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 1000")
	assert.Equal(t, message.SuccessMsg.Message, msg.Message, "Message must be null")
	assert.Equal(t, 3, len(couriers), "Count of courier must be 3")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")

}
