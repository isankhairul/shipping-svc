package entity

import (
	"go-klikdokter/app/model/base"
)

// swagger:model ChannelCourierService
type ChannelCourierService struct {
	base.BaseIDModel

	// Id of the Courier
	// in: integer
	CourierID uint64 `gorm:"not null" json:"courier_id"`

	// Id of the Courier
	// in: integer
	CourierServiceID uint64 `gorm:"not null" json:"courier_service_id"`

	// Id of the Channel
	// in: integer
	ChannelID uint64 `gorm:"not null" json:"channel_id"`

	// Status of the Courier
	// in: number
	PriceInternal float64 `gorm:"not null,type:decimal(18,4),default:0" json:"price_internal"`

	// Status of the Courier
	// in: integer
	Status int `gorm:"not null;default:1" json:"status"`

	Courier *Courier `json:"-" gorm:"foreignKey:courier_id"`

	Channel *Channel `json:"-" gorm:"foreignKey:channel_id"`

	CourierService *CourierService `json:"-" gorm:"foreignKey:courier_service_id"`
}
