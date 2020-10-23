package marshal

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []interface{}{
		1,
		1.5,
		true,
		"hello",
		[]int{1, 2, 3, 4, 5},
		[]float32{1.1, 2.2, 3.3},
		[]bool{true, false, false},
		[]string{"hello", "world", "!"},
		[][]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
		},
		[][]string{
			{"Alice", "Bob", "Carol"},
			{"Donald", "Emily", "Franklin"},
			{"George", "Helen", "Ivan"},
		},
	}

	for _, testCase := range testCases {
		// from encoding/json.Marshal
		b, err := json.Marshal(testCase)
		if err != nil {
			t.Error(err)
		}
		expected := string(b)

		// mine
		valField := reflect.ValueOf(testCase)
		mine := parse(valField)

		if mine != expected {
			t.Error("Not right")
		}
	}

}
