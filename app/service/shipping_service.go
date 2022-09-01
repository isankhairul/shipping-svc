package service

import (
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
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

	//Populate courier
	//courier map[id]code
	courier := make(map[uint64]string)
	var courierList []entity.Courier
	for _, v := range courierServices {
		_, ok := courier[v.CourierID]
		if !ok {
			courier[v.CourierID] = v.CourierCode
			courierList = append(courierList, entity.Courier{
				BaseIDModel: base.BaseIDModel{
					ID:  v.CourierID,
					UID: v.CourierUID,
				},
				Code: v.CourierCode,
			})
		}
	}

	price := s.GetAllCourierPrice(courierList, &input)
	resp = response.ToGetShippingRateResponseList(courierServices, price)

	return resp, message.SuccessMsg
}

// Get Price Data By Courier
func (s *shippingServiceImpl) GetAllCourierPrice(courier []entity.Courier, input *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	logger := log.With(s.logger, "ShippingService", "getAllCourierPrice")
	var resp = &response.ShippingRateCommonResponse{
		Rate:    make(map[string]response.ShippingRateData),
		Summary: map[string]response.ShippingRateSummary{},
	}

	for _, v := range courier {
		var (
			couriePrice *response.ShippingRateCommonResponse
			err         error
			log         string
		)

		// try to get price data from cache
		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%f",
			v.Code,
			input.Origin.PostalCode,
			input.Destination.PostalCode,
			input.Origin.Latitude,
			input.Origin.Longitude,
			input.Destination.Latitude,
			input.Destination.Longitude,
			input.TotalHeight,
		)

		_ = s.redis.GetJsonStruct(key, &couriePrice)
		// if cache doesn't exist
		if couriePrice == nil {
			switch v.Code {
			case shipping_provider.ShipperCode:
				couriePrice, err = s.shipper.GetShippingRate(&v.ID, input)
			default:
				continue
			}

			if err != nil {
				_ = level.Error(logger).Log(log, err.Error())
				continue
			}

			// save price to redis cache
			s.redis.SetJsonStruct(key, couriePrice, viper.GetInt("cache.redis.expired-in-minute.shipping-rate"))
		}

		if couriePrice != nil {
			resp.Add(couriePrice.Rate)
		}
	}

	return resp
}
