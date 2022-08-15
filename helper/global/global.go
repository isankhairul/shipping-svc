package global

import (
	"fmt"
	"html"
	"reflect"
	"strings"
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

func AddLike(column string, value []string) string {
	var condition string
	for _, v := range value {
		condition += fmt.Sprintf(" LOWER(%s) ILIKE '%%%s%%' OR", column, v)
	}
	return strings.TrimRight(condition, " OR")
}
