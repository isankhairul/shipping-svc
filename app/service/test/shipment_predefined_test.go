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

var baseshipmentPredefinedRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var shipmentPredefinedRepository = &repository_mock.ShipmentPredefinedMock{Mock: mock.Mock{}}
var shipmentPredefinedService = service.NewShipmentPredefinedService(logger, baseshipmentPredefinedRepository, shipmentPredefinedRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
}

func TestUpdateShipmentPredefined(t *testing.T) {
	item := entity.ShippmentPredefined{
		Type:   "type",
		Code:   "code",
		Title:  "title",
		Note:   "note",
		Status: 0,
	}
	req := request.UpdateShipmentPredefinedRequest{
		Uid:    "UCMvWngocMqKbaC3AWQBF",
		Type:   "type 1",
		Code:   "code 1",
		Title:  "title 1",
		Note:   "note 1",
		Status: 0,
	}

	shipmentPredefinedRepository.Mock.On("GetShipmentPredefinedByUid", mock.Anything).Return(item)
	shipmentPredefinedRepository.Mock.On("UpdateShipmentPredefined", mock.Anything).Return(item)
	result, _ := shipmentPredefinedService.UpdateShipmentPredefined(req)
	assert.NotNil(t, result)
	assert.Equal(t, "type", result.Type, "Type is type")
	assert.Equal(t, "code", result.Code, "Code is code")
	assert.Equal(t, "title", result.Title, "Title is title")
	assert.Equal(t, "note", result.Note, "Note is note")
}

func TestGetAll(t *testing.T) {
	req := request.ListShipmentPredefinedRequest{
		Page:  1,
		Limit: 10,
	}

	items := []*entity.ShippmentPredefined{
		{
			Type:   "type",
			Code:   "code",
			Title:  "title",
			Note:   "note",
			Status: 0,
		},
		{
			Type:   "type 1",
			Code:   "code 1",
			Title:  "title 1",
			Note:   "note 1",
			Status: 0,
		},
		{
			Type:   "type 2",
			Code:   "code 2",
			Title:  "title 2",
			Note:   "note 2",
			Status: 1,
		},
	}
	paginationResult := base.Pagination{
		Records:   120,
		Limit:     10,
		Page:      1,
		TotalPage: 12,
	}

	shipmentPredefinedRepository.Mock.On("GetAll", 10, 1, "").Return(items, &paginationResult)
	predefines, pagination, msg := shipmentPredefinedService.GetAll(req)

	assert.Equal(t, message.SuccessMsg.Code, msg.Code, "Code must be 201000")
	assert.Equal(t, 3, len(predefines), "Count of predefines must be 3")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")
}
