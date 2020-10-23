package marshal

import (
	"fmt"
	"reflect"
)

// Parse receives reflect.Value of a field,
// returns the converted string
func parse(valField reflect.Value) string {
	var ans string
	switch valField.Kind() {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14:
		ans = fmt.Sprint(valField)
	case reflect.String:
		ans = `"` + valField.String() + `"`
	case reflect.Array, reflect.Slice:
		ans = `[`
		for i := 0; i < valField.Len(); i++ {
			ans += parse(valField.Index(i))
			if i != valField.Len()-1 {
				ans += ","
			}
		}
		ans += `]`
	}
	return ans
}
