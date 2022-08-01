package entity

import "go-klikdokter/app/model/base"

// swagger:model ShippmentPredefined
type ShippmentPredefined struct {
	base.BaseIDModel
	// Name of the ShippmentPredefined
	// in: string
	Type string `gorm:"type:varchar(50);size:50;not null" json:"type"`

	// Code of the ShippmentPredefined
	// in: string
	Code string `gorm:"type:varchar(50);size:50;not null" json:"code"`

	// Type of the Courier
	// in: string
	Title string `gorm:"type:varchar(100);size:100;not null" json:"title"`

	// Note of the Courier
	// in: string
	Note string `gorm:"type:varchar(100);size:100;null" json:"note"`

	// Status of the ShippmentPredefined
	// in: integer
	Status int32 `gorm:"type:int;not null;default:1" json:"status"`
}
