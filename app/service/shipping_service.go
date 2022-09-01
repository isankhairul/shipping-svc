package service

import (
	"fmt"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/http_helper/shipping_provider"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/cache"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
)

type ShippingService interface {
	GetShippingRate(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message)
}

type shippingServiceImpl struct {
	logger              log.Logger
	baseRepo            repository.BaseRepository
	channelRepo         repository.ChannelRepository
	courierServiceRepo  repository.CourierServiceRepository
	courierCoverageCode repository.CourierCoverageCodeRepository
	shipper             shipping_provider.Shipper
	redis               cache.RedisCache
}

func NewShippingService(
	l log.Logger,
	br repository.BaseRepository,
	chrp repository.ChannelRepository,
	csrp repository.CourierServiceRepository,
	cccrp repository.CourierCoverageCodeRepository,
	sh shipping_provider.Shipper,
	rc cache.RedisCache,
) ShippingService {
	return &shippingServiceImpl{
		l, br, chrp, csrp, cccrp, sh, rc,
	}
}

// swagger:route POST /shipping/shipping-rate Shipping ShippingRate
// Get Shipping Rate
//
// responses:
//  200: ShippingRate
func (s *shippingServiceImpl) GetShippingRate(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetShippingRate")
	var resp []response.GetShippingRateResponse

	if len(input.ChannelCourierService) == 0 {
		return nil, message.ErrCourierServiceIsRequired
	}

	//Find Channel By UID
	channel, err := s.channelRepo.FindByUid(&input.ChannelUID)
	msg := message.SuccessMsg

	if err != nil {
		_ = level.Error(logger).Log("s.channelRepo.FindByUid", err.Error())
		return nil, message.ErrChannelNotFound
	}

	if channel == nil {
		return nil, message.ErrChannelNotFound
	}

	var courierServiceUIDs []string
	for _, v := range input.ChannelCourierService {
		courierServiceUIDs = append(courierServiceUIDs, v.CourierServiceUID)
	}

	//Find Courier Servies By Channel UID and Courier Servies UID Slice
	courierServices, err := s.courierServiceRepo.FindCourierServiceByChannelAndUIDs(input.ChannelUID, courierServiceUIDs)

	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierServiceByChannelAndUIDs", err.Error())
		return nil, message.ErrCourierServiceNotFound
	}

	if len(courierServices) == 0 {
		return nil, message.ErrCourierServiceNotFound
	}

	price := s.getPrice(&input, courierServices)
	resp = response.ToGetShippingRateResponseList(courierServices, price)

	return resp, msg
}

// Find Price Data
func (s *shippingServiceImpl) getPrice(input *request.GetShippingRateRequest, courierService []entity.ChannelCourierServiceForShippingRate) *response.ShippingRateCommonResponse {
	resp := response.ShippingRateCommonResponse{}
	resp.Rate = make(map[string]response.ShippingRateData)
	resp.Summary = make(map[string]response.ShippingRateSummary)

	for _, v := range courierService {

		originData, destinationData, msg := s.getOriginAndDestination(v.CourierID, input.Origin, input.Destination)

		courierShippingCode := global.CourierShippingCodeKey(v.CourierCode, v.ShippingCode)
		if msg.Code != message.SuccessMsg.Code {
			resp.Rate[courierShippingCode] = response.ShippingRateData{
				AvailableCode: 400,
				Error: response.GetShippingRateError{
					Message: msg.Message,
				},
			}
			continue
		}

		resp.Rate[courierShippingCode] = s.getSinglePrice(originData, destinationData, v.CourierCode, v.ShippingCode, input)
		resp.SummaryPerShippingType(v.ShippingTypeCode, resp.Rate[courierShippingCode].TotalPrice, v.EtdMax, v.EtdMin)
	}

	return &resp
}

// Get Origin and Destination data form database
func (s *shippingServiceImpl) getOriginAndDestination(courierID uint64, origin, destination request.AreaDetailPayload) (*entity.CourierCoverageCode, *entity.CourierCoverageCode, message.Message) {
	logger := log.With(s.logger, "ShippingService", "getOriginAndDestination")

	//Find Origin Area Data
	originResponse, err := s.courierCoverageCode.FindByCountryCodeAndPostalCode(courierID, origin.CountryCode, origin.PostalCode)

	if err != nil {
		_ = level.Error(logger).Log("s.courierCoverageCode.FindByCountryCodeAndPostalCode", err.Error())
		return nil, nil, message.ErrOriginNotFound
	}

	if originResponse == nil {
		return nil, nil, message.ErrOriginNotFound
	}

	//Find Destination Area Data
	destinationResponse, err := s.courierCoverageCode.FindByCountryCodeAndPostalCode(courierID, destination.CountryCode, destination.PostalCode)

	if err != nil {
		_ = level.Error(logger).Log("s.courierCoverageCode.FindByCountryCodeAndPostalCode", err.Error())
		return nil, nil, message.ErrDestinationNotFound
	}

	if destinationResponse == nil {
		return nil, nil, message.ErrDestinationNotFound

	}

	return originResponse, destinationResponse, message.SuccessMsg
}

// Get Third Party Price Data
func (s *shippingServiceImpl) getSinglePrice(origin, destination *entity.CourierCoverageCode, courierCode, shippingCode string, input *request.GetShippingRateRequest) response.ShippingRateData {

	var resp *response.ShippingRateCommonResponse
	var err error

	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%f",
		courierCode,
		input.Origin.PostalCode,
		input.Destination.PostalCode,
		input.Origin.Latitude,
		input.Origin.Longitude,
		input.Destination.Latitude,
		input.Destination.Longitude,
		input.TotalHeight,
	)

	// try to get data from cache
	err = s.redis.GetJsonStruct(key, &resp)
	courierShippingCode := global.CourierShippingCodeKey(courierCode, shippingCode)
	if resp != nil {
		return resp.FindShippingCode(courierShippingCode)
	}

	// get data from shipping provider if cache doesn't exist
	resp = &response.ShippingRateCommonResponse{
		Rate:    make(map[string]response.ShippingRateData),
		Summary: make(map[string]response.ShippingRateSummary),
	}
	switch courierCode {
	case shipping_provider.ShipperCode:
		resp, err = s.shipper.GetShippingRate(origin, destination, input)
	default:
		resp.Rate[courierShippingCode] = response.ShippingRateData{
			AvailableCode: 200,
			Error:         response.GetShippingRateError{Message: "shipping price not found"},
		}
	}

	// save data to cache
	if err == nil {
		s.redis.SetJsonStruct(key, resp, viper.GetInt("cache.redis.expired-in-minute.shipping-rate"))
	}

	return resp.FindShippingCode(courierShippingCode)
}
