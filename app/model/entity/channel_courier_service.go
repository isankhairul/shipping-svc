package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/pkg/util/datatype"
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

type ChannelCourierServiceForShippingRate struct {
	ShippingTypeCode        string         `gorm:"column:shipping_type_code"`
	ShippingTypeName        string         `gorm:"column:shipping_type_name"`
	ShippingTypeDescription string         `gorm:"column:shipping_type_description"`
	CourierID               uint64         `gorm:"column:courier_id"`
	CourierUID              string         `gorm:"column:courier_uid"`
	CourierCode             string         `gorm:"column:courier_code"`
	CourierName             string         `gorm:"column:courier_name"`
	CourierTypeCode         string         `gorm:"column:courier_type_code"`
	CourierTypeName         string         `gorm:"column:courier_type_name"`
	CourierServiceUID       string         `gorm:"column:courier_service_uid"`
	ShippingCode            string         `gorm:"column:shipping_code"`
	ShippingName            string         `gorm:"column:shipping_name"`
	ShippingDescription     string         `gorm:"column:shipping_description"`
	Logo                    datatype.JSONB `gorm:"column:logo"`
	EtdMin                  float64        `gorm:"column:etd_min"`
	EtdMax                  float64        `gorm:"column:etd_max"`
	Price                   float64        `gorm:"column:price"`
	UseInsurance            bool           `gorm:"column:use_insurance"`
	InsuranceFee            float64        `gorm:"column:insurance_fee"`
}
