package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ChannelService interface {
	CreateChannel(input request.SaveChannelRequest) (*entity.Channel, message.Message)
	GetList(input request.ChannelListRequest) ([]entity.Channel, *base.Pagination, message.Message)
	UpdateChannel(uid string, input request.SaveChannelRequest) message.Message
	GetChannel(uid string) (*entity.Channel, message.Message)
	DeleteChannel(uid string) message.Message
}

type ChannelServiceImpl struct {
	logger      log.Logger
	baseRepo    repository.BaseRepository
	channelRepo repository.ChannelRepository
}

func NewChannelService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.ChannelRepository,
) ChannelService {
	return &ChannelServiceImpl{lg, br, pr}
}

// swagger:route POST /channel/channel-app Channel ManagedChannelRequest
// Manage Channel
//
// responses:
//  401: SuccessResponse
//  201: SuccessResponse
func (s *ChannelServiceImpl) CreateChannel(input request.SaveChannelRequest) (*entity.Channel, message.Message) {
	logger := log.With(s.logger, "ChannelService", "CreateChannel")
	//Check exits `channel_code`
	//Set default value
	defaultLimit := 10
	defaultPage := 1
	defaultSort := ""
	filter := map[string]interface{}{
		"channel_code": input.ChannelCode,
	}

	result, _, err := s.channelRepo.FindByParams(defaultLimit, defaultPage, defaultSort, filter)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.FailedMsg
	}
	if result != nil {
		_ = level.Error(logger).Log(err)
		return nil, message.ErrDataChannelExists
	}

	s.baseRepo.BeginTx()
	//Set request to entity
	Channel := entity.Channel{
		ChannelName: input.ChannelName,
		ChannelCode: input.ChannelCode,
		Description: input.Description,
		Logo:        input.Logo,
		Status:      input.Status,
	}

	resultInsert, err := s.channelRepo.CreateChannel(&Channel)
	if err != nil {
		_ = level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.ErrDB
	}
	s.baseRepo.CommitTx()

	return resultInsert, message.SuccessMsg
}

// swagger:route GET /channel/channel-app/{uid} Get-channel Channel
// Get Channel
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
		"channel_code": input.ChannelCode,
		"status":       input.Status,
	}

	result, pagination, err := s.channelRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
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

// swagger:route PUT /channel/channel-app/{uid} UpdateChannelRequest
// Update Channel
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ChannelServiceImpl) UpdateChannel(uid string, input request.SaveChannelRequest) message.Message {
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

	data := map[string]interface{}{
		"channel_name": input.ChannelName,
		"channel_code": input.ChannelCode,
		"description":  input.Description,
		"logo":         input.Logo,
		"status":       input.Status,
	}

	err = s.channelRepo.Update(uid, data)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.FailedMsg
}

// swagger:route DELETE /channel/channel-app/{uid} channel-delete byParamDelete
// Delete Channel
//
// responses:
//  401: SuccessResponse
//  200: SuccessResponse
func (s *ChannelServiceImpl) DeleteChannel(uid string) message.Message {
	logger := log.With(s.logger, "ChannelService", "DeleteChannel")

	_, err := s.channelRepo.FindByUid(&uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	err = s.channelRepo.Delete(uid)
	if err != nil {
		_ = level.Error(logger).Log(err)
		return message.FailedMsg
	}

	return message.SuccessMsg
}
