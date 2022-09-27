package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/pkg/util/datatype"
)

type ShippingCourierStatus struct {
	base.BaseIDModel
	ShippingStatusID uint64          `gorm:"type:bigint;not null"`
	CourierID        uint64          `gorm:"type:bigint;not null"`
	StatusCode       string          `gorm:"type:varchar(100);size:100;not null"`
	StatusCourier    datatype.JSONB  `gorm:"type:jsonb;not null"`
	ShippingStatus   *ShippingStatus `gorm:"foreignKey:shipping_status_id"`
	Courier          *Courier        `gorm:"foreignKey:courier_id"`
}

func (ShippingCourierStatus) TableName() string {
	return "shipping_courier_status"
}
