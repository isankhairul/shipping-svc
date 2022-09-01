package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ChannelCourierServiceService interface {
	CreateChannelCourierService(input request.SaveChannelCourierServiceRequest) (*response.ChannelCourierServiceDetail, message.Message)
	ListChannelCouriersService(input request.ChannelCourierServiceListRequest) ([]response.ChannelCourierServiceItem, *base.Pagination, message.Message)
	GetChannelCourierService(uid string) (*response.ChannelCourierServiceDetail, message.Message)
	UpdateChannelCourierService(input request.UpdateChannelCourierServiceRequest) (*response.ChannelCourierServiceDetail, message.Message)
	DeleteChannelCourierService(uid string) message.Message
}

type channelCourierServiceServiceImpl struct {
	logger                 log.Logger
	baseRepo               repository.BaseRepository
	channelCouriers        repository.ChannelCourierRepository
	channelCourierServices repository.ChannelCourierServiceRepository
	courierServices        repository.CourierServiceRepository
}

func NewChannelCourierServiceService(
	lg log.Logger,
	br repository.BaseRepository,
	ccr repository.ChannelCourierRepository,
	channelCourierServices repository.ChannelCourierServiceRepository,
	courierServices repository.CourierServiceRepository,
) ChannelCourierServiceService {
	return &channelCourierServiceServiceImpl{lg, br, ccr, channelCourierServices, courierServices}
}

// swagger:route POST /channel/channel-courier-service/ Channel-Courier-Service CreateChannelCourierService
// Create Channel Courier Service
//
// responses:
//  200: ChannelCourierServiceDetailResponse
func (s *channelCourierServiceServiceImpl) CreateChannelCourierService(input request.SaveChannelCourierServiceRequest) (*response.ChannelCourierServiceDetail, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "CreateChannelCourierService")

	courierService, err := s.courierServices.FindByUid(&input.CourierServiceUID)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	if courierService == nil {
		return nil, message.ErrCourierServiceNotFound
	}

	channelCourier, err := s.channelCouriers.GetChannelCourierByUID(input.ChannelCourierUID)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	if channelCourier == nil {
		return nil, message.ErrChannelCourierNotFound
	}

	if channelCourier.Courier.UID != courierService.Courier.UID {
		return nil, message.ErrCourierServiceNotMatch
	}

	channelCourierService, err := s.channelCourierServices.GetChannelCourierService(channelCourier.ID, courierService.ID)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	if channelCourierService != nil {
		return response.NewChannelCourierServiceDetail(*channelCourierService), message.ErrDataExists
	}

	channelCourierService = &entity.ChannelCourierService{
		ChannelCourierID: channelCourier.ID,
		CourierServiceID: courierService.ID,
		PriceInternal:    input.PriceInternal,
		Status:           &input.Status,
	}

	result, err := s.channelCourierServices.CreateChannelCourierService(channelCourierService)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	result.ChannelCourier = channelCourier
	result.CourierService = courierService
	return response.NewChannelCourierServiceDetail(*channelCourierService), message.SuccessMsg
}

// swagger:route GET /channel/channel-courier-service/ Channel-Courier-Service GetChannelCourierServiceList
// Get List of Channel Courier Service
//
// responses:
//  200: ChannelCourierServiceList
func (s *channelCourierServiceServiceImpl) ListChannelCouriersService(input request.ChannelCourierServiceListRequest) ([]response.ChannelCourierServiceItem, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "ListChannelCouriersService")
	filters := map[string]interface{}{
		"status":        input.Filters.Status,
		"shipping_name": input.Filters.ShippingName,
		"shipping_code": input.Filters.ShippingCode,
		"shipping_type": input.Filters.ShippingType,
		"courier_name":  input.Filters.CourierName,
		"channel_name":  input.Filters.ChannelName,
		"courier_uid":   input.Filters.CourierUID,
	}

	result, paging, err := s.channelCourierServices.FindByParams(input.Limit, input.Page, input.Sort, filters)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, nil, message.ErrDB
	}

	return response.NewChannelCourierServiceList(result), paging, message.SuccessMsg
}

// swagger:route GET /channel/channel-courier-service/{uid} Channel-Courier-Service GetChannelCourierServiceByUID
// Get Detail of Channel Courier Service
//
// responses:
//  200: ChannelCourierServiceDetailResponse
func (s *channelCourierServiceServiceImpl) GetChannelCourierService(uid string) (*response.ChannelCourierServiceDetail, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "GetChannelCourierService")

	result, err := s.channelCourierServices.GetChannelCourierServiceByUID(uid)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return response.NewChannelCourierServiceDetail(*result), message.SuccessMsg
}

// swagger:route PUT /channel/channel-courier-service/{uid} Channel-Courier-Service UpdateChannelCourierService
// Update a channel courier by uid
//
// responses:
//  200: ChannelCourierServiceDetailResponse
func (s *channelCourierServiceServiceImpl) UpdateChannelCourierService(input request.UpdateChannelCourierServiceRequest) (*response.ChannelCourierServiceDetail, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "UpdateChannelCourierService")

	data := map[string]interface{}{
		"status":         input.Body.Status,
		"price_internal": input.Body.PriceInternal,
	}

	result, err := s.channelCourierServices.GetChannelCourierServiceByUID(input.UID)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrChannelCourierNotFound
	}

	err = s.channelCourierServices.UpdateChannelCourierService(input.UID, data)
	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return nil, message.ErrDB
	}

	return s.GetChannelCourierService(input.UID)
}

// swagger:route DELETE /channel/channel-courier-service/{uid} Channel-Courier-Service DeleteChannelCourierServiceByUID
// Delete Courier Service
//
// responses:
//  200: SuccessResponse
func (s *channelCourierServiceServiceImpl) DeleteChannelCourierService(uid string) message.Message {
	logger := log.With(s.logger, "ChannelCourierService", "DeleteChannelCourierService")

	result, err := s.channelCourierServices.GetChannelCourierServiceByUID(uid)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return message.ErrDB
	}

	if result == nil {
		return message.ErrNoData
	}

	err = s.channelCourierServices.DeleteChannelCourierServiceByID(result.ID)
	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return message.ErrDB
	}

	return message.SuccessMsg
}
