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

type ChannelService interface {
	GetList(input request.ChannelListRequest) ([]entity.Channel, *base.Pagination, message.Message)
	GetListStatus(input request.GetChannelCourierStatusRequest) ([]response.GetChannelCourierStatusResponseItem, *base.Pagination, message.Message)
	GetChannel(uid string) (*entity.Channel, message.Message)
	CreateChannel(input request.SaveChannelRequest) (*entity.Channel, message.Message)
	UpdateChannel(input request.UpdateChannelRequest) message.Message
	DeleteChannel(uid string) message.Message
}

type ChannelServiceImpl struct {
	logger                log.Logger
	baseRepo              repository.BaseRepository
	channelRepo           repository.ChannelRepository
	shippingCourierStatus repository.ShippingCourierStatusRepository
}

func NewChannelService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.ChannelRepository,
	scs repository.ShippingCourierStatusRepository,
) ChannelService {
	return &ChannelServiceImpl{lg, br, pr, scs}
}

// swagger:route GET /channel/channel-app Channel-Apps Channels
// List of Channel Apps
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *ChannelServiceImpl) GetList(input request.ChannelListRequest) ([]entity.Channel, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ChannelService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	filter := map[string]interface{}{
		"channel_code": input.Filters.ChannelCode,
		"channel_name": input.Filters.ChannelName,
		"status":       input.Filters.Status,
	}
	result, pagination, err := s.channelRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, nil, message.FailedMsg
	}
	if len(result) == 0 {
		_ = level.Warn(logger).Log(message.ErrNoData)
		return nil, nil, message.ErrNoData
	}

	return result, pagination, message.SuccessMsg
}

// swagger:route GET /channel/channel-app/{uid} Channel-Apps ChannelRequestGetByUid
// Get Detail of Channel Apps
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *ChannelServiceImpl) GetChannel(uid string) (*entity.Channel, message.Message) {
	logger := log.With(s.logger, "ChannelService", "GetChannel")
	result, err := s.channelRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	if result == nil {
		return nil, message.ErrNoData
	}

	return result, message.SuccessMsg
}

// swagger:route POST /channel/channel-app Channel-Apps SaveChannelRequest
// Add Channel App
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *ChannelServiceImpl) CreateChannel(input request.SaveChannelRequest) (*entity.Channel, message.Message) {
	logger := log.With(s.logger, "ChannelService", "CreateChannel")
	//Check exits `channel_code`
	//Set default value
	uid := "0" //set UID 0 as a default to check duplicate channel code

	isExists, err := s.channelRepo.CheckExistsByUIdChannelCode(uid, input.ChannelCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}
	if isExists {
		_ = level.Error(logger).Log(message.ErrNoData)
		return nil, message.ErrDataChannelExists
	}

	//Set request to entity
	channel := entity.Channel{
		ChannelName: input.ChannelName,
		ChannelCode: input.ChannelCode,
		Description: input.Description,
		Logo:        input.Logo,
		Status:      1, //Default
		ImageUID:    input.ImageUID,
		ImagePath:   input.ImagePath,
	}
	if input.Status != nil {
		channel.Status = *input.Status
	}
	resultInsert, err := s.channelRepo.CreateChannel(&channel)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDB
	}

	return resultInsert, message.SuccessMsg
}

// swagger:route PUT /channel/channel-app/{uid} Channel-Apps UpdateChannelRequest
// Update Channel App by uid
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ChannelServiceImpl) UpdateChannel(input request.UpdateChannelRequest) message.Message {
	uid := input.Uid
	logger := log.With(s.logger, "ChannelService", "UpdateChannel")
	channel, err := s.channelRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}
	if channel == nil {
		_ = level.Error(logger).Log(message.ErrNoData)
		return message.FailedMsg
	}

	//Check exists channel_code
	isExists, err := s.channelRepo.CheckExistsByUIdChannelCode(uid, input.ChannelCode)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}
	if isExists {
		_ = level.Error(logger).Log(message.ErrNoData)
		return message.ErrDataChannelExists
	}

	data := map[string]interface{}{
		"channel_name": input.ChannelName,
		"channel_code": input.ChannelCode,
		"description":  input.Description,
		"logo":         input.Logo,
		"status":       input.Status,
		"image_uid":    input.ImageUID,
		"image_path":   input.ImagePath,
	}

	err = s.channelRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}

// swagger:route DELETE /channel/channel-app/{uid} Channel-Apps ChannelRequestDeleteByUid
// Delete Channel Apps by uid
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ChannelServiceImpl) DeleteChannel(uid string) message.Message {
	logger := log.With(s.logger, "ChannelService", "DeleteChannel")
	channel, err := s.channelRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	if channel == nil {
		return message.ErrChannelNotFound
	}

	hasChild := s.channelRepo.IsChannelHasChild(channel.ID)

	if hasChild.ChannelCourier {
		return message.ErrChannelHasCourierAssigned
	}

	if hasChild.ShippingStatus {
		return message.ErrChannelHasChildShippingStatus
	}

	err = s.channelRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}

// swagger:route GET /channel/channel-status-courier-status Channel-Apps GetChannelCourierStatus
// Get Channel Courier Status List
//
// responses:
//  200: GetChannelCourierStatusResponse
func (s *ChannelServiceImpl) GetListStatus(input request.GetChannelCourierStatusRequest) ([]response.GetChannelCourierStatusResponseItem, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ChannelService", "GetListStatus")

	filters := map[string]interface{}{
		"channel_name": input.Filters.ChannelName,
		"courier_name": input.Filters.CourierName,
		"status_code":  input.Filters.StatusCode,
	}

	result, paging, err := s.shippingCourierStatus.FindByParams(input.Limit, input.Page, input.Sort, filters)

	if err != nil {
		_ = level.Error(logger).Log(err.Error())
		return []response.GetChannelCourierStatusResponseItem{}, nil, message.FailedMsg
	}

	if len(result) == 0 {
		return []response.GetChannelCourierStatusResponseItem{}, nil, message.ErrNoData
	}

	return response.NewGetChannelCourierStatusResponse(result), paging, message.SuccessMsg
}
