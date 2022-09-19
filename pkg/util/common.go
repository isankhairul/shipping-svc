package util

import "reflect"

func IsSliceAndNotEmpty(input interface{}) bool {
	if reflect.ValueOf(input).Kind() == reflect.Slice {
		return reflect.ValueOf(input).Len() > 0
	}
	return false
}

// Replace empty string with default value
func ReplaceEmptyString(str string, defaultValue string) string {
	if len(str) == 0 {
		return defaultValue
	}

	return str
}
