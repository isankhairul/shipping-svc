package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/helper/message"
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

	// readonly field
	ShippingTypeName string `gorm:"-:migration;->"`
}

type ChannelCourierServiceForShippingRate struct {
	ShippingTypeCode            string         `gorm:"column:shipping_type_code"`
	ShippingTypeName            string         `gorm:"column:shipping_type_name"`
	ShippingTypeDescription     string         `gorm:"column:shipping_type_description"`
	CourierID                   uint64         `gorm:"column:courier_id"`
	CourierUID                  string         `gorm:"column:courier_uid"`
	CourierCode                 string         `gorm:"column:courier_code"`
	CourierName                 string         `gorm:"column:courier_name"`
	CourierTypeCode             string         `gorm:"column:courier_type_code"`
	CourierTypeName             string         `gorm:"column:courier_type_name"`
	CourierServiceUID           string         `gorm:"column:courier_service_uid"`
	ShippingCode                string         `gorm:"column:shipping_code"`
	ShippingName                string         `gorm:"column:shipping_name"`
	ShippingDescription         string         `gorm:"column:shipping_description"`
	Logo                        datatype.JSONB `gorm:"column:logo"`
	EtdMin                      float64        `gorm:"column:etd_min"`
	EtdMax                      float64        `gorm:"column:etd_max"`
	Price                       float64        `gorm:"column:price"`
	UseInsurance                bool           `gorm:"column:use_insurance"`
	InsuranceFee                float64        `gorm:"column:insurance_fee"`
	MaxWeight                   float64        `gorm:"column:max_weight"`
	CourierStatus               int32          `gorm:"column:courier_status"`
	CourierServiceStatus        int32          `gorm:"column:courier_service_status"`
	ChannelCourierStatus        int32          `gorm:"column:channel_courier_status"`
	ChannelCourierServiceStatus int32          `gorm:"column:channel_courier_service_status"`
	HidePurpose                 int32          `gorm:"column:hide_purpose"`
	PrescriptionAllowed         int32          `gorm:"column:prescription_allowed"`
}

func (c *ChannelCourierServiceForShippingRate) Validate(finalWeight *float64, prescription_allowed *bool) message.Message {

	if msg, isValid := c.IsValidCourier(); !isValid {
		return msg
	}

	if msg, isValid := c.IsValidCourierService(finalWeight, prescription_allowed); !isValid {
		return msg
	}

	return message.SuccessMsg
}

// validate courier data
func (c *ChannelCourierServiceForShippingRate) IsValidCourier() (message.Message, bool) {

	if c.CourierStatus != 1 {
		return message.CourierNotActiveMsg, false
	}

	if c.HidePurpose != 0 {
		return message.CourierHiddenInPurposeMsg, false
	}

	if c.ChannelCourierStatus != 1 {
		return message.ChannelCourierNotActiveMsg, false
	}

	return message.SuccessMsg, true
}

// validate courier service data
func (c *ChannelCourierServiceForShippingRate) IsValidCourierService(finalWeight *float64, prescription_allowed *bool) (message.Message, bool) {

	if c.CourierServiceStatus != 1 {
		return message.CourierServiceNotActiveMsg, false
	}

	if c.ChannelCourierServiceStatus != 1 {
		return message.ChannelCourierServiceNotActiveMsg, false
	}

	if c.MaxWeight > 0 && c.MaxWeight < *finalWeight {
		return message.WeightExceedsmsg, false
	}

	if *prescription_allowed && c.PrescriptionAllowed != 1 {
		return message.PrescriptionNotAllowedMsg, false
	}

	return message.SuccessMsg, true
}
