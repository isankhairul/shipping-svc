package global

import (
	"html"
	"reflect"
)

func HtmlEscape(req interface{}) {
	value := reflect.ValueOf(req).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}

		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))
	}
}
