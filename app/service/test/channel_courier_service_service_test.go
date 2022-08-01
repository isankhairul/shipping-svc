package test

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
}

func TestCreateChannelCourierServiceSuccess(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierServiceRequest{
		CourierServiceUID: "CourierServiceUID",
		ChannelCourierUID: "ChannelCourierUID",
		PriceInternal:     1000,
		Status:            1,
	}

	channelCourier := &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cc_1"}, Courier: &entity.Courier{CourierName: "Courier_1"}, Channel: &entity.Channel{ChannelName: "Channel_1"}}
	courierService := &entity.CourierService{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cs_1"}, Courier: &entity.Courier{CourierName: "Courier_1"}}
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(*courierService)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(channelCourier)

	//return nil when creating
	channelCourierServiceRepo.Mock.On("GetChannelCourierService", mock.Anything, mock.Anything).Return(nil)

	result, msg := channelCourierServiceService.CreateChannelCourierService(input)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
	assert.Equal(t, "Courier_1", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "Channel_1", result.ChannelName, "CourierName must be test name")
}

func TestCreateChannelCourierServiceCourierNotMatch(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierServiceRequest{
		CourierServiceUID: "CourierServiceUID",
		ChannelCourierUID: "ChannelCourierUID",
		PriceInternal:     1000,
		Status:            1,
	}

	channelCourier := &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cc_1"}, Courier: &entity.Courier{BaseIDModel: base.BaseIDModel{ID: 1, UID: "1"}, CourierName: "Courier_1"}, Channel: &entity.Channel{ChannelName: "Channel_1"}}
	courierService := &entity.CourierService{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cs_1"}, Courier: &entity.Courier{BaseIDModel: base.BaseIDModel{ID: 2, UID: "2"}, CourierName: "Courier_2"}}
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(*courierService)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(channelCourier)

	//return nil when creating
	channelCourierServiceRepo.Mock.On("GetChannelCourierService", mock.Anything, mock.Anything).Return(nil)

	result, msg := channelCourierServiceService.CreateChannelCourierService(input)

	assert.Nil(t, result)
	assert.Equal(t, 34402, msg.Code, "Status code should be 34402")
}

func TestCreateChannelCourierServiceFailedChannelCourierServiceExist(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierServiceRequest{
		CourierServiceUID: "CourierServiceUID",
		ChannelCourierUID: "ChannelCourierUID",
		PriceInternal:     1000,
		Status:            1,
	}

	channelCourier := &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cc_1"}, Courier: &entity.Courier{CourierName: "Courier_1"}, Channel: &entity.Channel{ChannelName: "Channel_1"}}
	courierService := &entity.CourierService{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cs_1"}, Courier: &entity.Courier{CourierName: "Courier_1"}}
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(*courierService)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(channelCourier)

	cc := &entity.ChannelCourierService{
		BaseIDModel:      base.BaseIDModel{UID: "123"},
		ChannelCourierID: channelCourier.ID,
		CourierServiceID: courierService.ID,
		PriceInternal:    input.PriceInternal,
		Status:           &input.Status,
		ChannelCourier:   channelCourier,
		CourierService:   courierService,
	}
	//return nil when creating
	channelCourierServiceRepo.Mock.On("GetChannelCourierService", mock.Anything, mock.Anything).Return(cc)

	result, msg := channelCourierServiceService.CreateChannelCourierService(input)

	assert.NotNil(t, result)
	assert.Equal(t, 34001, msg.Code, "Status code should be 34001")
	assert.Equal(t, "Courier_1", result.CourierName, "CourierName must be test name")
	assert.Equal(t, "Channel_1", result.ChannelName, "CourierName must be test name")
}

func TestCreateChannelCourierServiceFailedChannelCourierNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierServiceRequest{
		CourierServiceUID: "CourierServiceUID",
		ChannelCourierUID: "ChannelCourierUID",
		PriceInternal:     1000,
		Status:            1,
	}

	courierService := &entity.CourierService{BaseIDModel: base.BaseIDModel{ID: 1, UID: "cs_1"}}
	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(*courierService)
	channelCourierRepo.Mock.On("GetChannelCourierByUID", mock.Anything).Return(nil)

	result, msg := channelCourierServiceService.CreateChannelCourierService(input)

	assert.Nil(t, result)
	assert.Equal(t, 34402, msg.Code, "Status code should be 34402")
}

func TestCreateChannelCourierServiceFailedCourierServiceNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.SaveChannelCourierServiceRequest{
		CourierServiceUID: "CourierServiceUID",
		ChannelCourierUID: "ChannelCourierUID",
		PriceInternal:     1000,
		Status:            1,
	}

	courierServiceRepo.Mock.On("FindByUid", mock.Anything).Return(nil)

	result, msg := channelCourierServiceService.CreateChannelCourierService(input)

	assert.Nil(t, result)
	assert.Equal(t, 34101, msg.Code, "Status code should be 34101")
}

func TestListChannelCouriersServiceSuccess(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.ChannelCourierServiceListRequest{
		ChannelName:  []string{},
		CourierName:  []string{},
		Status:       []int{},
		ShippingName: []string{},
		ShippingCode: []string{},
		ShippingType: []string{},
	}

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	status1 := int32(1)
	ccs := entity.ChannelCourierService{
		BaseIDModel:    base.BaseIDModel{UID: "a"},
		ChannelCourier: &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: "b"}, Courier: courier, Channel: channel},
		CourierService: &entity.CourierService{BaseIDModel: base.BaseIDModel{UID: "c"}},
		Status:         &status1,
	}
	channelCourierServiceRepo.Mock.On("FindByParams", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.ChannelCourierService{ccs}, &base.Pagination{}, nil)
	result, _, msg := channelCourierServiceService.ListChannelCouriersService(input)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
}

func TestListChannelCouriersServiceNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.ChannelCourierServiceListRequest{
		ChannelName:  []string{},
		CourierName:  []string{},
		Status:       []int{},
		ShippingName: []string{},
		ShippingCode: []string{},
		ShippingType: []string{},
	}

	channelCourierServiceRepo.Mock.On("FindByParams", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]entity.ChannelCourierService{}, &base.Pagination{}, nil)
	result, _, msg := channelCourierServiceService.ListChannelCouriersService(input)

	assert.Nil(t, result)
	assert.Equal(t, 34005, msg.Code, "Status code should be 34005")
}

func TestGetChannelCourierServiceSuccess(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	status1 := int32(1)
	ccs := entity.ChannelCourierService{
		BaseIDModel:    base.BaseIDModel{UID: "a"},
		ChannelCourier: &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: "b"}, Courier: courier, Channel: channel},
		CourierService: &entity.CourierService{BaseIDModel: base.BaseIDModel{UID: "c"}},
		Status:         &status1,
	}
	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(&ccs, nil)
	result, msg := channelCourierServiceService.GetChannelCourierService("uid")

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
}

func TestGetChannelCourierServiceNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(nil, nil)
	result, msg := channelCourierServiceService.GetChannelCourierService("uid")

	assert.Nil(t, result)
	assert.Equal(t, 34005, msg.Code, "Status code should be 34005")
}

func TestUpdateChannelCourierServiceSuccess(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierServiceRequest{
		UID: "a",
		Body: request.UpdateChannelCourierService{
			PriceInternal: 1,
			Status:        1,
		},
	}

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	status1 := int32(1)
	ccs := entity.ChannelCourierService{
		BaseIDModel:    base.BaseIDModel{UID: "a"},
		ChannelCourier: &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: "b"}, Courier: courier, Channel: channel},
		CourierService: &entity.CourierService{BaseIDModel: base.BaseIDModel{UID: "c"}},
		Status:         &status1,
	}
	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(&ccs, nil)
	channelCourierServiceRepo.Mock.On("UpdateChannelCourierService", mock.Anything, mock.Anything).Return(nil)
	result, msg := channelCourierServiceService.UpdateChannelCourierService(input)

	assert.NotNil(t, result)
	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
}

func TestUpdateGetChannelCourierServiceByUIDNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierServiceRequest{
		UID: "a",
		Body: request.UpdateChannelCourierService{
			PriceInternal: 1,
			Status:        1,
		},
	}

	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(nil, nil)
	result, msg := channelCourierServiceService.UpdateChannelCourierService(input)

	assert.Nil(t, result)
	assert.Equal(t, 34402, msg.Code, "Status code should be 34402")
}

func TestDeleteChannelCourierServiceSuccess(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierServiceRequest{
		UID: "a",
		Body: request.UpdateChannelCourierService{
			PriceInternal: 1,
			Status:        1,
		},
	}

	courier := &entity.Courier{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	channel := &entity.Channel{
		BaseIDModel: base.BaseIDModel{UID: "a"},
	}

	status1 := int32(1)
	ccs := entity.ChannelCourierService{
		BaseIDModel:    base.BaseIDModel{UID: "a"},
		ChannelCourier: &entity.ChannelCourier{BaseIDModel: base.BaseIDModel{UID: "b"}, Courier: courier, Channel: channel},
		CourierService: &entity.CourierService{BaseIDModel: base.BaseIDModel{UID: "c"}},
		Status:         &status1,
	}
	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(&ccs, nil)
	channelCourierServiceRepo.Mock.On("DeleteChannelCourierServiceByID", mock.Anything).Return(nil)
	msg := channelCourierServiceService.DeleteChannelCourierService(input.UID)

	assert.Equal(t, 201000, msg.Code, "Status code should be 200")
}

func TestDeleteChannelCourierServiceFailedNotFound(t *testing.T) {
	var channelCourierRepo = &repository_mock.ChannelCourierRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceRepo = &repository_mock.ChannelCourierServiceRepositoryMock{Mock: mock.Mock{}}
	var courierServiceRepo = &repository_mock.CourierServiceRepositoryMock{Mock: mock.Mock{}}
	var channelCourierServiceService = service.NewChannelCourierServiceService(logger, baseRepository, channelCourierRepo, channelCourierServiceRepo, courierServiceRepo)

	input := request.UpdateChannelCourierServiceRequest{
		UID: "a",
		Body: request.UpdateChannelCourierService{
			PriceInternal: 1,
			Status:        1,
		},
	}

	channelCourierServiceRepo.Mock.On("GetChannelCourierServiceByUID", mock.Anything).Return(nil, nil)
	msg := channelCourierServiceService.DeleteChannelCourierService(input.UID)

	assert.Equal(t, 34005, msg.Code, "Status code should be 34005")
}
