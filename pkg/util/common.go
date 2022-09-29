package util

import (
	"reflect"
)

func IsSliceAndNotEmpty(input interface{}) bool {
	if reflect.ValueOf(input).Kind() == reflect.Slice {
		return reflect.ValueOf(input).Len() > 0
	}
	return false
}

func IsNilOrEmpty(input interface{}) bool {

	if input == nil {
		return true
	}

	if reflect.ValueOf(input).IsZero() {
		return true
	}

	switch reflect.ValueOf(input).Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map, reflect.String:
		return reflect.ValueOf(input).Len() == 0
	}

	return true
}

// Replace empty string with default value
func ReplaceEmptyString(str string, defaultValue string) string {
	if len(str) == 0 {
		return defaultValue
	}

	return str
}
