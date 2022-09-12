package entity

import "go-klikdokter/app/model/base"

type OrderShippingHistory struct {
	base.BaseIDModel
	OrderShippingID  uint64 `gorm:"type:bigint;not null"`
	ShippingStatusID uint64 `gorm:"type:bigint;not null"`
	StatusCode       string `gorm:"type:varchar(50);not null"`
	Note             string `gorm:"type:varchar(255);null"`

	ShippingCourierStatus *ShippingCourierStatus `gorm:"foreignKey:shipping_status_id"`
}

func (OrderShippingHistory) TableName() string {
	return "order_shipping_history"
}
