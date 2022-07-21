package datatype

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Time string

var timeFormat = "15:04:05-07"

func TimeNow(loc ...*time.Location) Time {
	location := time.UTC

	if loc != nil {
		location = loc[0]
	}

	return Time(time.Now().In(location).Format(timeFormat))
}

func (Time) GormDataType() string {
	return "time"
}

func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "TIME"
	case "postgres":
		return "TIME(0) WITH TIME ZONE"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

func (t Time) String() string {
	return string(t)
}
