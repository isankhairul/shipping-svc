package entity

import (
	"go-klikdokter/app/model/base"
)

// swagger:model ChannelCourierService
type ChannelCourierService struct {
	base.BaseIDModel

	// Id of the Channel Courier
	// in: integer
	ChannelCourierID uint64 `gorm:"type:bigint;not null" json:"channel_courier_id"`

	// Id of the Courier
	// in: integer
	CourierServiceID uint64 `gorm:"type:bigint;not null" json:"courier_service_id"`

	// Status of the Courier
	// in: number
	PriceInternal float64 `gorm:"type:numeric;not null,type:decimal(18,4),default:0" json:"price_internal"`

	// Status of the Courier
	// in: integer
	Status *int32 `gorm:"type:int;not null;default:1" json:"status"`

	ChannelCourier *ChannelCourier `json:"-" gorm:"foreignKey:channel_courier_id"`

	CourierService *CourierService `json:"-" gorm:"foreignKey:courier_service_id"`
}
