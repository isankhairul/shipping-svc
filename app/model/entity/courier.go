package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/pkg/util/datatype"
)

// swagger:model Courier
type Courier struct {
	base.BaseIDModel
	// Name of the Courier
	// in: string
	CourierName string `gorm:"not null" json:"courier_name"`

	// Code of the Courier
	// in: string
	Code string `gorm:"not null;unique:true" json:"code"`

	// Type of the Courier
	// in: string
	CourierType string `gorm:"not null" json:"courier_type"`

	// Logo of the Courier
	// in: string
	Logo string `gorm:"null;type:varchar(500)" json:"logo"`

	// Hide purpose of the Courier
	// in: integer
	HidePurpose int `gorm:"not null;default:0" json:"hide_purpose"`

	// Courier Api Integration of the Courier. Need to set column becase snake-case cannot understand convetnion.
	// in: integer
	CourierApiIntegration int `gorm:"not null;default:1;column:api_integration" json:"courier_api_intergration"`

	// Geo Coodinate of the Courier
	// in: string
	UseGeocoodinate int `gorm:"not null;default:0" json:"use_geocoodinate"`

	// Provide Airwaybill of the Courier
	// in: integer
	ProvideAirwaybill int `gorm:"not null;default:0" json:"provide_airwaybill"`

	// Status of the Courier
	// in: integer
	Status *int `gorm:"not null;default:1" json:"status"`

	CourierCoverageCode []*CourierCoverageCode `gorm:"foreignKey:courier_uid" json:"-"`

	CourierServices []*CourierService `json:"-" gorm:"foreignKey:courier_id"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`
}
