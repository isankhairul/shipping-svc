package entity

import (
	"go-klikdokter/app/model/base"
)

// swagger:model ChannelCourier
type ChannelCourier struct {
	base.BaseIDModel

	// Id of the Courier
	// in: integer
	CourierID uint64 `gorm:"not null" json:"courier_id"`

	// Id of the Channel
	// in: integer
	ChannelID uint64 `gorm:"not null" json:"channel_id"`

	// PrioritySort of the ChannelCourier
	// in: int
	PrioritySort int `gorm:"not null;default:1" json:"priority_sort"`

	// Hide purpose of the ChannelCourier
	// in: integer
	HidePurpose int `gorm:"not null;default:0" json:"hide_purpose"`

	// Status of the ChannelCourier
	// in: integer
	Status *int `gorm:"not null;default:1" json:"status"`

	Courier *Courier `json:"-" gorm:"foreignKey:courier_id"`

	Channel *Channel `json:"-" gorm:"foreignKey:channel_id"`

	ChannelCourierServices []*ChannelCourierService `json:"-"`
}

// swagger: model ChannelCourierServiceDTO
type ChannelCourierServiceDTO struct {
	CourierServiceUID string  `json:"courier_service_uid"`
	CourierUID        string  `json:"courier_uid"`
	ChannelUID        string  `json:"channel_uid"`
	ShippingType      string  `json:"shipping_type"`
	ShippingName      string  `json:"shipping_name"`
	ShippingCode      string  `json:"shipping_code"`
	PriceInternal     float64 `json:"price_internal"`
	Status            int     `json:"status"`
}

func ToChannelCourierServiceDTO(channelCourierService *ChannelCourierService, courierService *CourierService) *ChannelCourierServiceDTO {
	ret := &ChannelCourierServiceDTO{
		CourierServiceUID: courierService.UID,
		ShippingType:      courierService.ShippingType,
		ShippingName:      courierService.ShippingName,
		ShippingCode:      courierService.ShippingCode,
		PriceInternal:     channelCourierService.PriceInternal,
		Status:            *channelCourierService.Status,
	}
	return ret
}

// swagger: model ChannelCourierDTO
type ChannelCourierDTO struct {
	Uid                string                      `json:"uid"`
	ChannelName        string                      `json:"channel_name"`
	ChannelCode        string                      `json:"channel_code"`
	CourierName        string                      `json:"courier_name"`
	PrioritySort       int                         `json:"priority_sort"`
	HidePurpose        int                         `json:"hide_purpose"`
	Status             int                         `json:"status"`
	CourierServiceDTOs []*ChannelCourierServiceDTO `json:"-"`
}

func ToChannelCourierDTO(cur *ChannelCourier) *ChannelCourierDTO {
	ret := &ChannelCourierDTO{
		Uid:          cur.UID,
		PrioritySort: cur.PrioritySort,
		HidePurpose:  cur.HidePurpose,
		Status:       *cur.Status,
	}

	if cur.Channel != nil {
		ret.ChannelName = cur.Channel.ChannelName
		ret.ChannelCode = cur.Channel.ChannelCode
	}

	if cur.Courier != nil {
		ret.CourierName = cur.Courier.CourierName
	}

	return ret
}
