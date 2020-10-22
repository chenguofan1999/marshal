package marshal

import (
	"errors"
	"fmt"
	"reflect"
)

// JsonMarshal receives an arbitrary struct and return Json format byte stream
func JsonMarshal(v interface{}) ([]byte, error) {
	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)

	vKind := vVal.Kind()
	if vKind != reflect.Struct {
		return nil, errors.New("Expect a struct")
	}

	//字段数
	num := vVal.NumField()
	// fmt.Printf("有 %d 个字段\n", num)

	s := "{"
	for i := 0; i < num; i++ {
		//fmt.Println(vType.Field(i).Name, " ",vVal.Field(i)," ", vVal.Field(i).Kind(), "| ", vType.Field(i).Tag.Get("myTag"))
		s += "\""
		valField := vVal.Field(i)
		typeField := vType.Field(i)

		tagName, ok := typeField.Tag.Lookup("myTag")
		if ok {
			s += tagName
		} else {
			s += typeField.Name
		}

		s += "\":"

		switch valField.Kind() {
		case reflect.String:
			s += "\""
			s += valField.String()
			s += "\""
		case reflect.Slice, reflect.Array:
			s += "["

			// Get element type
			elemKind := typeField.Type.Elem().Kind()

			for i := 0; i < valField.Len(); i++ {

				if elemKind == reflect.String {
					s += "\""
				}

				s += fmt.Sprint(valField.Index(i))

				if elemKind == reflect.String {
					s += "\""
				}

				if i != valField.Len()-1 {
					s += ","
				}
			}

			s += "]"

		default:
			s += fmt.Sprint(valField)
		}

		if i != num-1 {
			s += ","
		}
	}
	s += "}"
	return []byte(s), nil
}
