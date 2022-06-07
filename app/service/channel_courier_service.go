package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ChannelCourierService interface {
	CreateChannelCourier(input request.SaveChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message)
	ListChannelCouriers(input request.ChannelCourierListRequest) ([]*entity.ChannelCourierDTO, *base.Pagination, message.Message)
	GetChannelCourier(uid string) (*entity.ChannelCourierDTO, message.Message)
	UpdateChannelCourier(input request.UpdateChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message)
	DeleteChannelCourier(uid string) message.Message
}

type ChannelCourierServiceImpl struct {
	logger                 log.Logger
	baseRepo               repository.BaseRepository
	channelCouriers        repository.ChannelCourierRepository
	channelCourierServices repository.ChannelCourierServiceRepository
	courierServices        repository.CourierServiceRepository
}

func NewChannelCourierService(
	lg log.Logger,
	br repository.BaseRepository,
	ccr repository.ChannelCourierRepository,
	channelCourierServices repository.ChannelCourierServiceRepository,
	courierServices repository.CourierServiceRepository,
) ChannelCourierService {
	return &ChannelCourierServiceImpl{lg, br, ccr, channelCourierServices, courierServices}
}

// swagger:route POST /channel/channel-courier/ Channel SaveChannelCourierRequest
// Assign Courier to Channel
//
// responses:
//  401: errorResponse
//  200: ChannelCourier
func (s *ChannelCourierServiceImpl) CreateChannelCourier(input request.SaveChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message) {
	ret, msg := s.createChannelCourierInTx(input)
	return ret, msg
}

func (s *ChannelCourierServiceImpl) createChannelCourierInTx(input request.SaveChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "CreateChannelCourier")

	courier, notFoundCourier := s.channelCouriers.FindCourierByUID(input.CourierUID)
	if notFoundCourier != nil {
		return nil, message.ErrCourierNotFound
	}
	channel, notFoundChannel := s.channelCouriers.FindChannelByUID(input.ChannelUID)
	if notFoundChannel != nil {
		return nil, message.ErrChannelNotFound
	}

	cc, _ := s.channelCouriers.GetChannelCourierByIds(channel.ID, courier.ID)
	if cc != nil {
		return entity.ToChannelCourierDTO(cc), message.ErrChannelCourierFound
	}

	cc = &entity.ChannelCourier{
		CourierID:    courier.ID,
		ChannelID:    channel.ID,
		PrioritySort: input.PrioritySort,
		HidePurpose:  input.HidePurpose,
		Status:       input.Status,
	}
	for _, courierServiceUID := range input.CourierServiceUIDs {
		courierService, err := s.courierServices.FindByUid(&courierServiceUID.CourierServiceUid)
		if courierService == nil {
			return nil, message.ErrNoDataCourierService
		}
		if err != nil || courierService.Status == 0 {
			return nil, message.ErrCourierServiceHasInvalidStatus
		}
		_, err = s.channelCourierServices.CreateChannelCourierService(courier, channel, courierService, courierServiceUID.PriceInternal, courierServiceUID.Status)
		if err != nil {
			return nil, message.ErrChannelCourierServiceCreateFailed
		}
	}
	cc, err := s.channelCouriers.CreateChannelCourier(cc)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}
	cc, err = s.channelCouriers.GetChannelCourierByUID(cc.UID)
	if err != nil {
		_ = level.Error(logger).Log(err)
	}
	return entity.ToChannelCourierDTO(cc), message.SuccessMsg
}

// swagger:route GET /channel/channel-courier/{uid} Channel GetChannelCourierByUid
// Get Detail of Channel Courier
//
// responses:
//  200: ChannelCourierDTO
//  401: UnauthorizedResponse
// 	400: InvalidRequestDataResponse
//  500: InternalServerErrorResponse
func (s *ChannelCourierServiceImpl) GetChannelCourier(uid string) (*entity.ChannelCourierDTO, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "Get Detail of Channel Courier")
	cur, err := s.channelCouriers.GetChannelCourierByUID(uid)
	if cur == nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrChannelCourierNotFound
	}
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}
	return entity.ToChannelCourierDTO(cur), message.SuccessMsg
}

// swagger:route GET /channel/channel-courier Channel ChannelCourierListRequest
// List of Assignment Channel and Courier
//
// responses:
//  200: PaginationResponse
//  401: UnauthorizedResponse
// 	400: InvalidRequestDataResponse
//  500: InternalServerErrorResponse
func (s *ChannelCourierServiceImpl) ListChannelCouriers(input request.ChannelCourierListRequest) ([]*entity.ChannelCourierDTO, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "Get Detail of Channel Courier")

	filter := map[string]interface{}{
		"status":       input.Status,
		"courier_name": input.CourierName,
		"channel_name": input.ChannelName,
		"channel_code": input.ChannelCode,
	}

	result, pagination, err := s.channelCouriers.FindByPagination(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, nil, message.FailedMsg
	}

	if result == nil {
		_ = level.Warn(logger).Log(message.ErrNoData)
		return nil, nil, message.FailedMsg
	}
	items := make([]*entity.ChannelCourierDTO, len(result))
	for index, element := range result {
		items[index] = entity.ToChannelCourierDTO(element)
	}
	return items, pagination, message.SuccessMsg
}

// swagger:route PUT /channel/channel-courier/{uid} Channel UpdateChannelCourierRequest
// Update a channel courier by uid
//
// responses:
//  401: errorResponse
//  201: ChannelCourier
func (s *ChannelCourierServiceImpl) UpdateChannelCourier(input request.UpdateChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message) {
	ret, msg := s.updateChannelCourierInTx(input)
	return ret, msg
}

func (s *ChannelCourierServiceImpl) updateChannelCourierInTx(input request.UpdateChannelCourierRequest) (*entity.ChannelCourierDTO, message.Message) {
	logger := log.With(s.logger, "ChannelCourierService", "UpdateChannelCourier")
	data := map[string]interface{}{
		"hide_purpose":  input.HidePurpose,
		"status":        input.Status,
		"priority_sort": input.PrioritySort,
	}
	cur, err := s.channelCouriers.GetChannelCourierByUID(input.Uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrChannelCourierNotFound
	}
	result := s.channelCouriers.UpdateChannelCourier(input.Uid, data)
	if result != nil {
		_ = level.Error(logger).Log(message.ErrNoData)
		return nil, message.ErrChannelCourierNotFound
	}

	for _, courierServiceUID := range input.CourierServiceUIDs {
		courierService, err := s.courierServices.FindByUid(&courierServiceUID.CourierServiceUid)
		if courierService == nil {
			return nil, message.ErrNoDataCourierService
		}
		if err != nil || courierService.Status == 0 {
			return nil, message.ErrCourierServiceHasInvalidStatus
		}
		_, err = s.channelCourierServices.CreateChannelCourierService(cur.Courier, cur.Channel, courierService,
			float64(courierServiceUID.PriceInternal), courierServiceUID.Status)
		if err != nil {
			_ = level.Error(logger).Log(err)
			return nil, message.ErrChannelCourierServiceCreateFailed
		}
	}
	//deleting old channelCourierServices
	if cur.ChannelCourierServices != nil {
		inputUIDs := mapInputUIDS(input.CourierServiceUIDs)
		for _, ccs := range cur.ChannelCourierServices {
			if !contains(&ccs.CourierService.UID, inputUIDs) {
				err := s.channelCourierServices.DeleteChannelCourierServiceByID(ccs.ID)
				if err != nil {
					_ = level.Error(logger).Log(err)
				}
			}
		}
	}
	cur, err = s.channelCouriers.GetChannelCourierByUID(input.Uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
	}
	return entity.ToChannelCourierDTO(cur), message.SuccessMsg
}

// swagger:route DELETE /channel/channel-courier/{uid} Channel DeleteChannelCourierByUid
// Delete Channel Courier
//
// responses:
//  200: SuccessResponse
//  401: UnauthorizedResponse
// 	400: InvalidRequestDataResponse
//  500: InternalServerErrorResponse
func (s *ChannelCourierServiceImpl) DeleteChannelCourier(uid string) message.Message {
	channelCourier, err := s.channelCouriers.GetChannelCourierByUID(uid)
	if err != nil {
		return message.ErrChannelCourierNotFound
	}
	err = s.channelCourierServices.DeleteChannelCourierServicesByChannelID(channelCourier.ChannelID, channelCourier.CourierID)
	if err != nil {
		return message.ErrUnableToDeleteChannelCourier
	}
	err = s.channelCouriers.DeleteChannelCourierByID(channelCourier.ID)
	if err != nil {
		return message.ErrUnableToDeleteChannelCourier
	}
	return message.SuccessMsg
}

func contains(cur *string, items []*string) bool {
	for _, value := range items {
		if strings.Compare(*cur, *value) == 0 {
			return true
		}
	}

	return false
}

func mapInputUIDS(courierServiceUIDs []*request.CourierServiceDTO) []*string {
	items := []*string{}
	for _, value := range courierServiceUIDs {
		items = append(items, &value.CourierServiceUid)
	}
	return items
}
