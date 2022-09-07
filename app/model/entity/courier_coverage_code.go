package entity

import "go-klikdokter/app/model/base"

// swagger:model CourierCoverageCode
type CourierCoverageCode struct {
	base.BaseIDModel

	// Relation with CourierID
	// in: Courier
	// require: true
	CourierID uint64 `gorm:"type:bigint;not null" json:"-"`

	// Courier UID of the Courier
	// in: string
	// require: false
	// example: ggkjhsdf6668885555
	CourierUID string `gorm:"-:all" json:"courier_uid"` //this field will be ignored by gorm

	CourierName string `gorm:"-" json:"courier_name"`

	// Country code of the Courier Coverage Code
	// in: string
	// require: true
	// example: ID
	CountryCode string `gorm:"type:varchar(20) not null" json:"country_code"`

	// Postal code of the Courier Coverage Code
	// in: string
	// require: true
	// example: 151338
	PostalCode string `gorm:"type:varchar(20) not null" json:"postal_code"`

	// Subdistrict of the Courier Coverage Code
	// in: string
	// require: true
	// example: 151338
	Subdistrict string `gorm:"type:varchar(50) default('') not null" json:"subdistrict"`

	// Description of the Courier Coverage Code
	// in: string
	// example: PAGEDANGAN
	Description string `gorm:"type:varchar(100)" json:"description"`

	// Code 1 of the Courier Coverage Code
	// in: string
	// example: CKG011
	Code1 string `gorm:"type:varchar(50)" json:"code1"`

	// Code 2 of the Courier Coverage Code
	// in: string
	// example: CKG012
	Code2 string `gorm:"type:varchar(50)" json:"code2"`

	// Code 3 of the Courier Coverage Code
	// in: string
	Code3 string `gorm:"type:varchar(50)" json:"code3"`

	// Code 4 of the Courier Coverage Code
	// in: string
	Code4 string `gorm:"type:varchar(50)" json:"code4"`

	// Code 5 of the Courier Coverage Code
	// in: string
	Code5 string `gorm:"type:varchar(50)" json:"code5"`

	// Code 6 of the Courier Coverage Code
	// in: string
	Code6 string `gorm:"type:varchar(50)" json:"code6"`

	// Status of the Courier
	// in: integer
	Status *int32 `gorm:"type:int;not null;default:1" json:"status"`

	Courier *Courier `json:"-" gorm:"foreignKey:courier_id"`
}
