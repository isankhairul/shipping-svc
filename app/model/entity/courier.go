package entity

import "go-klikdokter/app/model/base"

// swagger:model Courier
type Courier struct {
	base.BaseIDModel
	// Name of the Courier
	// in: string
	CourierName string `gorm:"not null" json:"courier_name"`

	// Code of the Courier
	// in: string
	Code string `gorm:"unique,not null" json:"code"`

	// Type of the Courier
	// in: string
	CourierType string `gorm:"not null" json:"courier_type"`

	// Logo of the Courier
	// in: string
	Logo string `json:"logo"`

	// Hide purpose of the Courier
	// in: integer
	HidePurpose int `gorm:"not null;default:0" json:"hide_purpose"`

	// Courier Api Integration of the Courier
	// in: integer
	CourierApiIntegration int `gorm:"not null;default:1" json:"courier_api_intergration"`

	// Geo Coodinate of the Courier
	// in: string
	UseGeocoodinate int `gorm:"not null;default:0" json:"use_geocoodinate"`

	// Provide Airwaybill of the Courier
	// in: integer
	ProvideAirwaybill int `gorm:"not null;default:0" json:"provide_airwaybill"`

	// Status of the Courier
	// in: integer
	Status int `gorm:"not null;default:1" json:"status"`
}