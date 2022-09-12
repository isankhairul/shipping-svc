package entity

import "go-klikdokter/app/model/base"

type OrderShippingItem struct {
	base.BaseIDModel
	OrderShippingID uint64  `gorm:"type:bigint;not null"`
	ItemName        string  `gorm:"type:varchar(100);not null"`
	ProductUID      string  `gorm:"type:varchar(50);null"`
	Price           float64 `gorm:"type:numeric;not null"`
	Quantity        int     `gorm:"type:integer;not null"`
	UOM             string  `gorm:"type:varchar(50);not null"`
	Prescription    int     `gorm:"type:integer;not null"`
	TotalPrice      float64 `gorm:"type:numeric;null"`
	Weight          float64 `gorm:"type:numeric;null"`
	Volume          float64 `gorm:"type:numeric;null"`
}

func (OrderShippingItem) TableName() string {
	return "order_shipping_item"
}
