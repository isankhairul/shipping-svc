package response

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/pkg/util/datatype"
)

// swagger:response GetChannelCourierStatusResponse
type GetChannelCourierStatusResponse struct {
	// in: body
	Body []GetChannelCourierStatusResponseItem `json:"body"`
}

// swagger:model GetChannelCourierStatusResponse
type GetChannelCourierStatusResponseItem struct {
	// Channel code
	ChannelCode string `json:"channel_code"`

	// Channel name
	ChannelName string `json:"channel_name"`

	// Courier name
	CourierName string `json:"courier_name"`

	// Status Code
	StatusCode string `json:"status_code"`

	// Status Title
	StatusTitle string `json:"status_title"`

	// Courier Status
	CourierStatus datatype.JSONB `json:"courier_status"`
}

func NewGetChannelCourierStatusResponseItem(input entity.ShippingCourierStatus) GetChannelCourierStatusResponseItem {
	resp := GetChannelCourierStatusResponseItem{
		CourierStatus: input.StatusCourier,
	}

	if input.ShippingStatus != nil {

		if input.ShippingStatus.Channel != nil {
			resp.ChannelName = input.ShippingStatus.Channel.ChannelName
			resp.ChannelCode = input.ShippingStatus.Channel.ChannelCode
		}

		resp.StatusCode = input.ShippingStatus.StatusCode
		resp.StatusTitle = input.ShippingStatus.StatusName
	}

	if input.Courier != nil {
		resp.CourierName = input.Courier.CourierName
	}
	return resp
}

func NewGetChannelCourierStatusResponse(input []entity.ShippingCourierStatus) []GetChannelCourierStatusResponseItem {
	resp := []GetChannelCourierStatusResponseItem{}
	for _, v := range input {
		resp = append(resp, NewGetChannelCourierStatusResponseItem(v))
	}
	return resp
}
