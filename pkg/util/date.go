package util

import (
	"time"
)

var (
	LayoutDefault  = "2006-01-02 15:04:05"
	LayoutDateOnly = "2006-01-02"
	LayoutTimeOnly = "15:04"
	Loc, _         = time.LoadLocation("Asia/Jakarta")
)

func DateRangeValidation(dateStart, dateEnd string) (validDate, validRange bool) {
	validDate = true
	validRange = false

	dateStartT, err := time.ParseInLocation(LayoutDateOnly, dateStart, Loc)
	if err != nil {
		validDate = false
	}

	dateEndT, err := time.ParseInLocation(LayoutDateOnly, dateEnd, Loc)
	if err != nil {
		validDate = false
	}

	if dateEndT.After(dateStartT) {
		validRange = true
	}

	return validDate, validRange
}

func DateValidationYYYYMMDD(date string) (validDate bool) {
	validDate = true

	_, err := time.ParseInLocation(LayoutDateOnly, date, Loc)
	if err != nil {
		validDate = false
	}
	return validDate
}

func ConvertToDateTime(dateTime string) (time.Time, error) {
	DateTime, err := time.ParseInLocation(LayoutDefault, dateTime, Loc)
	if err != nil {
		return time.Time{}, err
	}

	return DateTime, err
}
