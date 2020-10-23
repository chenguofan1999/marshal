package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWithoutTags(t *testing.T) {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
		IsWarm bool
	}
	redGroup := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
		IsWarm: true,
	}

	// Mine
	b, err := JsonMarshal(redGroup)
	if err != nil {
		t.Error(err)
	}

	// Golang's
	expected, err := json.Marshal(redGroup)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(expected) {
		fmt.Println("Expected : ", expected)
		fmt.Println("Actual   : ", string(b))
		t.Error("Not right")
	}

}

func TestWithTags(t *testing.T) {
	type ColorGroup struct {
		ID     int      `myTag:"GroupID"`
		Name   string   `myTag:"GroupName"`
		Colors []string `myTag:"ColorNames"`
		IsWarm bool     `myTag:"Warm Color?"`
	}
	redGroup := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
		IsWarm: true,
	}
	b, err := JsonMarshal(redGroup)
	if err != nil {
		t.Error(err)
	}

	expected := `{"GroupID":1,"GroupName":"Reds","ColorNames":["Crimson","Red","Ruby","Maroon"],"Warm Color?":true}`

	if string(b) != expected {
		fmt.Println("Expected : ", expected)
		fmt.Println("Actual   : ", string(b))
		t.Error("Not right")
	}
}

func TestAllTypes(t *testing.T) {
	type TestGroup struct {
		IntVal   int
		Int8Val  int8
		Int16Val int16
		Int32Val int32
		Int64Val int64

		Float32Val float32
		Float64Val float64

		BoolVal     bool
		StringVal   string
		IntSlice    []int
		FloatSlice  []float32
		StringSlice []string
	}

	testGroup := TestGroup{
		IntVal:      200,
		Int8Val:     100,
		Int16Val:    300,
		Int32Val:    400,
		Int64Val:    500,
		Float32Val:  600.0,
		Float64Val:  700.0,
		BoolVal:     true,
		StringVal:   "testStringVal",
		IntSlice:    []int{1, 2, 3, 4, 5},
		FloatSlice:  []float32{6.0, 7.0, 8.0, 9.0},
		StringSlice: []string{"10", "20", "30", "40"},
	}

	// Mine
	b, err := JsonMarshal(testGroup)
	if err != nil {
		t.Error(err)
	}

	// Golang's
	expected, err := json.Marshal(testGroup)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(expected) {
		fmt.Println("Expected : ", string(expected))
		fmt.Println("Actual   : ", string(b))
		t.Error("Not right")
	}
}

func TestMoreData(t *testing.T) {
	type FruitBasket struct {
		Name    string
		Fruit   []string
		ID      int64
		Private string // An unexported field is not encoded.
	}

	basket := []FruitBasket{
		{
			Name:    "Standard",
			Fruit:   []string{"Apple", "Banana", "Orange"},
			ID:      999,
			Private: "Second-rate",
		}, {
			Name:    "Alice",
			Fruit:   []string{"strawberry", "watermelon", "Orange"},
			ID:      234,
			Private: "First-rate",
		}, {
			Name:    "Bob",
			Fruit:   []string{"blackberry", "melon", "cherry"},
			ID:      867,
			Private: "First-rate",
		}, {
			Name:    "Henry",
			Fruit:   []string{"pear", "pomegranate", "pitaya"},
			ID:      657,
			Private: "First-rate",
		},
	}

	for _, fruitBasket := range basket {

		expected, err := json.Marshal(fruitBasket)
		if err != nil {
			t.Error(err)
		}

		b, err := JsonMarshal(fruitBasket)
		if err != nil {
			t.Error(err)
		}

		if string(b) != string(expected) {
			fmt.Println("Expected : ", string(expected))
			fmt.Println("Actual   : ", string(b))
			t.Error("Not right")
		}
	}
}

func TestMultidimensionalSlice(t *testing.T) {
	type TestStruct struct {
		Ints    [][]int
		Floats  [][]float32
		Strings [][]string
	}
	testCase := TestStruct{
		Ints: [][]int{
			{0, 1, 2, 3},
			{4, 5, 6, 7},
			{8, 9, 10, 11},
		},
		Floats: [][]float32{
			{0.0, 1.0, 2.0},
			{3.0, 4.0, 5.0},
			{6.0, 7.0, 8.0},
		},
		Strings: [][]string{
			{"aa", "bb", "cc", "dd"},
			{"ee", "ff", "gg", "hh"},
		},
	}

	// Mine
	b, err := JsonMarshal(testCase)
	if err != nil {
		t.Error(err)
	}

	// Golang's
	expected, err := json.Marshal(testCase)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(expected) {
		fmt.Println("Expected : ", string(expected))
		fmt.Println("Actual   : ", string(b))
		t.Error("Not right")
	}
}
