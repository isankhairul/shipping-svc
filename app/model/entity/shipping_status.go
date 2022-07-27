package entity

import "go-klikdokter/app/model/base"

type ShippingStatus struct {
	base.BaseIDModel
	ChannelID   uint64   `gorm:"type:bigint;not null"`
	StatusCode  string   `gorm:"type:varchar;size:50;not null"`
	StatusName  string   `gorm:"type:varchar;size:100;not null"`
	Description string   `gorm:"type:varchar;size:100;not null"`
	Channel     *Channel `gorm:"foreignKey:channel_id"`
}

func (ShippingStatus) TableName() string {
	return "shipping_status"
}
