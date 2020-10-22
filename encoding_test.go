package marshal

import (
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
	b, err := JsonMarshal(redGroup)
	if err != nil {
		t.Error(err)
	}

	expected := `{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"],"IsWarm":true}`

	if string(b) != expected {
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
		Price  float32  `myTag:"PriceOfSTH"`
	}
	redGroup := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
		IsWarm: true,
		Price:  1.99,
	}
	b, err := JsonMarshal(redGroup)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Actual   : ", string(b))
}
