package marshal

import (
	"errors"
	"reflect"
)

// JsonMarshal receives an arbitrary struct and return Json format byte stream
func JsonMarshal(v interface{}) ([]byte, error) {
	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)

	// Check if input is a struct
	vKind := vVal.Kind()
	if vKind != reflect.Struct {
		return nil, errors.New("Expect a struct")
	}

	// Number of fields
	num := vVal.NumField()

	s := "{"
	for i := 0; i < num; i++ {
		s += `"`
		valField := vVal.Field(i)
		typeField := vType.Field(i)

		// Replace Field Name by myTag if exists
		tagName, ok := typeField.Tag.Lookup("myTag")
		if ok {
			s += tagName
		} else {
			s += typeField.Name
		}

		s += `":`

		s += parse(valField)

		if i != num-1 {
			s += ","
		}
	}
	s += "}"
	return []byte(s), nil
}
