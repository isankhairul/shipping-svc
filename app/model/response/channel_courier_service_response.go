package response

import "go-klikdokter/app/model/entity"

// swagger:response ChannelCourierServiceDetailResponse
type ChannelCourierServiceDetailResponse struct {
	// in: body
	Body ChannelCourierServiceDetail `json:"body"`
}

// swagger:model ChannelCourierServiceDetail
type ChannelCourierServiceDetail struct {
	UID               string  `json:"uid"`
	ChannelCourierUID string  `json:"channel_courier_uid"`
	CourierServiceUID string  `json:"courier_service_uid"`
	CourierName       string  `json:"courier_name"`
	ChannelName       string  `json:"channel_name"`
	ShippingName      string  `json:"shipping_name"`
	PriceInternal     float64 `json:"price_internal"`
	Status            int     `json:"status"`
}

func NewChannelCourierServiceDetail(input entity.ChannelCourierService) *ChannelCourierServiceDetail {
	response := ChannelCourierServiceDetail{
		UID:           input.UID,
		PriceInternal: input.PriceInternal,
		Status:        *input.Status,
	}

	if input.ChannelCourier != nil {
		response.ChannelCourierUID = input.ChannelCourier.UID
	}

	if input.ChannelCourier != nil && input.ChannelCourier.Channel != nil {
		response.ChannelName = input.ChannelCourier.Channel.ChannelName
	}

	if input.ChannelCourier != nil && input.ChannelCourier.Courier != nil {
		response.CourierName = input.ChannelCourier.Courier.CourierName
	}

	if input.CourierService != nil {
		response.CourierServiceUID = input.CourierService.UID
		response.ShippingName = input.CourierService.ShippingName
	}

	return &response
}

// swagger:response ChannelCourierServiceList
type ChannelCourierServiceList struct {
	// in: body
	Body []ChannelCourierServiceItem `json:"body"`
}

// swagger:model ChannelCourierServiceItem
type ChannelCourierServiceItem struct {
	UID          string `json:"uid"`
	ChannelName  string `json:"channel_name"`
	CourierName  string `json:"courier_name"`
	ShippingName string `json:"shipping_name"`
	ShippingCode string `json:"shipping_code"`
	ShippingType string `json:"shipping_type"`
	Status       int    `json:"status"`
}

func NewChannelCourierServiceItem(input entity.ChannelCourierService) *ChannelCourierServiceItem {
	response := ChannelCourierServiceItem{
		UID:    input.UID,
		Status: *input.Status,
	}

	if input.ChannelCourier != nil && input.ChannelCourier.Channel != nil {
		response.ChannelName = input.ChannelCourier.Channel.ChannelName
	}

	if input.ChannelCourier != nil && input.ChannelCourier.Courier != nil {
		response.CourierName = input.ChannelCourier.Courier.CourierName
	}

	if input.CourierService != nil {
		response.ShippingName = input.CourierService.ShippingName
		response.ShippingCode = input.CourierService.ShippingCode
		response.ShippingType = input.CourierService.ShippingType
	}

	return &response
}

func NewChannelCourierServiceList(input []entity.ChannelCourierService) []ChannelCourierServiceItem {
	response := []ChannelCourierServiceItem{}
	for _, v := range input {
		response = append(response, *NewChannelCourierServiceItem(v))
	}
	return response
}
