# 设计说明

本项目的主要目标是实现 `encoding/json` 里的 `Marshal` 函数，该函数接受一个 struct，返回该 struct 的 JSON 格式的字节流。


## `JsonMarshal` 函数

直接来看最终需要实现的 `JsonMarshal` 函数:

```go
func JsonMarshal(v interface{}) ([]byte, error) 
```
为了能接受各种各样的 struct , 该函数接受的参数实际上是一个空接口，因为任意一个结构都实现了空接口。但同时也带来了问题 —— 非结构的变量也实现了空接口，因此也可以被传入。因此对传入参数类型的检查是必要的：

```go
	vVal := reflect.ValueOf(v)

	// Check if input is a struct
	vKind := vVal.Kind()
	if vKind != reflect.Struct {
		return nil, errors.New("Expect a struct")
	}
```
通过 `reflect.ValueOf()` 得到传入参数的反射对象，并通过反射对象的 `Kind()` 方法得到参数的类型，如果不是一个 struct，则直接返回错误。

如果传入的参数是一个 struct，则需要遍历该结构体的各字段。当反射对象是 struct 类型时：
- 可用 `reflect.Value.NumField()` 方法获取其字段数，
- 可用 `reflect.Value.Field(i)` 获取其第 i 个字段的反射对象。

以下是遍历结构体 v 的各字段的基本方式：

```go
    vVal := reflect.ValueOf(v) // reflect.Value of v
    num := vVal.NumField()     // number of fields of v
    for i := 0; i < num; i++ {
        ...
        // reflect.Value of the ith field of v
		valField := vVal.Field(i)
        ...
    }
```

对于每一个字段，我们需要分别对字段名和字段值进行序列化。同时为了实现**字段标签**，需要检查有效标签的存在，如果存在，则需用其替代字段名。

获取字段名需要通过 `reflect.Type` , 而获取标签也需要通过该类型。  
以下代码展示了获取结构体 v 的各字段的字段名、标签名的基本方式：
```go
    vType := reflect.TypeOf(v) // reflect.Type of v
    num := vVal.NumField()     // number of fields of v
    for i := 0; i < num; i++ {
        ...

        // reflect.Type of the ith field of v
        typeField := vType.Field(i)

        // Field Name
        fieldName := typeField.Name

        // Tag Name
        tagName, ok := typeField.Tag.Lookup("myTag")

        ...
    }
```

然后是实现 `JsonMarshal` 的关键，即对字段数据的序列化。这里把这一模块独立出来，即 `parse()` 函数。

将上面的各部分按需求组合起来便是 `JsonMarshal` 函数的大体框架，如下是完整代码。

```go
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
```

## `parse` 函数

`parse` 函数接受结构体的一个字段的反射对象 `reflect.Value`，返回其序列化后的字符串。

需要对字段进行分类：
- 对于整数、实数、bool类型，直接返回其数值的字符串形式。
- 对于 string 类型，在前后添上引号后返回。
- 对于数组、切片类型，在首尾添上 '[' 和 ']'，递归地对每个元素应用 parse ，元素之间用 ',' 连接。（递归的设计同时自然地实现了对多维数组、切片的转换）


```go
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
```


**`reflect.Kind`：**

`reflect.Value.Kind()` 返回的 `reflect.Kind` 类型实质上是整数，例如对于 bool 类型，其对应的 `reflect.Kind` 既可以表示为 `reflect.Bool` ，也可以直接表示为 1 . 因此上面代码中的这一段

```go
case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14:
		ans = fmt.Sprint(valField)
```

表示的是整数、实数、bool类型的情况，其中 1 表示 bool 类型，2 - 11各自表示了不同格式的整数类型，13、14表示了两种浮点数类型。

# 单元测试

本包共包含两个模块，`JsonMarshal` 和 `parse` 两个函数。由于`JsonMarshal` 集成了 `parse` ，对 `JsonMarshal` 的单元测试可视为集成测试，同时也是本包的功能测试。

## 对 `parse` 的单元测试

测试数据包含各种具有代表性的数据类型：整数、实数、bool、string、这四者的一维切片、多维切片。

对于每个数据，同时用 `parse` 和 `encoding/json.Marshal` 函数生成结果，后者作为标准参考，对比两个结果，如果一致则测试通过。

```go
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
```

**测试结果:**

```
=== RUN   TestParse
--- PASS: TestParse (0.00s)
PASS
ok      command-line-arguments  0.003s
```

测试通过。


## 对 `JsonMarshal` 的测试

设计了多个测试函数和测例，分别对无标签普通情况、有标签情况、多维数组序列化进行了测试。其中对于无标签的测试，通过与 `encoding/json.Marshal` 函数的结果进行比对来验证；对于有标签的测试，由于我的标签格式与之不同，则自己写了预期结果与测试结果进行比对。

**无标签普通情况**
```go
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
```

**有字段标签**
```go
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

```

**更多数据的测试**
```go
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

```

**对多维数组字段的测试**
```go
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
```

**测试结果：**

```
$ go test -v                                       
=== RUN   TestParse
--- PASS: TestParse (0.00s)
=== RUN   TestWithoutTags
--- PASS: TestWithoutTags (0.00s)
=== RUN   TestWithTags
--- PASS: TestWithTags (0.00s)
=== RUN   TestAllTypes
--- PASS: TestAllTypes (0.00s)
=== RUN   TestMoreData
--- PASS: TestMoreData (0.00s)
=== RUN   TestMultidimensionalSlice
--- PASS: TestMultidimensionalSlice (0.00s)
PASS
ok      github.com/chenguofan1999/marshal       0.003s
```

测试通过。