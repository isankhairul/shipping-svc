package entity

import "go-klikdokter/app/model/base"

// swagger:model ShippmentPredefined
type ShippmentPredefined struct {
	base.BaseIDModel
	// Name of the ShippmentPredefined
	// in: string
	Type string `gorm:"not null,nvarchar(50)" json:"type"`

	// Code of the ShippmentPredefined
	// in: string
	Code string `gorm:"not null,nvarchar(50)" json:"code"`

	// Type of the Courier
	// in: string
	Title string `gorm:"not null,nvarchar(100)" json:"title"`

	// Note of the Courier
	// in: string
	Note string `gorm:"null,nvarchar(100)" json:"note"`

	// Status of the ShippmentPredefined
	// in: integer
	Status int `gorm:"not null;default:1" json:"status"`
}
