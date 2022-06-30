package entity

import "go-klikdokter/app/model/base"

// swagger:model Channel
type Channel struct {
	base.BaseIDModel
	// ChannelName of the Channel
	// in: string
	ChannelName string `gorm:"not null" json:"channel_name"`

	// ChannelCode of the Channel
	// in: string
	ChannelCode string `gorm:"unique,not null" json:"channel_code"`

	// Description of the Channel
	// in: string
	Description string `gorm:"not null" json:"description"`

	// Logo of the Channel
	// in: string
	Logo string `gorm:"not null" json:"logo"`

	// Status of the Channel
	// in: integer
	Status int `gorm:"not null" json:"status"`
}
