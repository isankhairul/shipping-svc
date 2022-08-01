package base

import (
	"time"

	"go-klikdokter/pkg/util"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type BaseIDModel struct {
	//gorm.Model
	// Id of as primary key
	// in: int64
	ID uint64 `gorm:"primary_key:auto_increment" json:"id"`

	// UID of the a model
	// in: int64
	UID string `gorm:"type:varchar(21);size:21;uniqueIndex" json:"uid"`

	IsDeleted bool      `json:"is_deleted" gorm:"type:boolean"`
	CreatedBy string    `json:"-" gorm:"type:varchar(100);size:100"`
	CreatedAt time.Time `json:"-" gorm:"type:timestamp"`
	UpdatedBy string    `json:"-" gorm:"type:varchar(100);size:100"`
	UpdatedAt time.Time `json:"-" gorm:"type:timestamp"`
}

func (base *BaseIDModel) BeforeCreate(tx *gorm.DB) error {
	uid, _ := gonanoid.New()
	tx.Statement.SetColumn("UID", uid)
	tx.Statement.SetColumn("IsDeleted", false)
	tx.Statement.SetColumn("CreatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	tx.Statement.SetColumn("UpdatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	return nil
}

func (base *BaseIDModel) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	return nil
}
