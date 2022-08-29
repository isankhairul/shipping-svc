package util

import "reflect"

func IsSliceAndNotEmpty(input interface{}) bool {
	if reflect.ValueOf(input).Kind() == reflect.Slice {
		return reflect.ValueOf(input).Len() > 0
	}
	return false
}
