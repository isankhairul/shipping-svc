package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/pkg/util/datatype"
)

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

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`
}