# marshal

[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/chenguofan1999/marshal)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/chenguofan1999/marshal)

```
                                   _               _ 
                                  | |             | |
  _ __ ___     __ _   _ __   ___  | |__     __ _  | |
 | '_ ` _ \   / _` | | '__| / __| | '_ \   / _` | | |
 | | | | | | | (_| | | |    \__ \ | | | | | (_| | | |
 |_| |_| |_|  \__,_| |_|    |___/ |_| |_|  \__,_| |_|
                                                     
                                                     
```

Golang encoding/json.Marshal的最简实现。



## 安装


```sh
go get github.com/chenguofan1999/marshal
```

**更新**

```sh
go get -u github.com/chenguofan1999/marshal
```

## 功能

提供 `JsonMarshal` 函数，用于将任意 struct 转换为 JSON 格式的字节流。

**支持的结构字段数据类型：**

- int
- int8
- int16
- int32
- int64
- uint16
- uint32
- uint64
- float32
- float64
- bool
- string
- 以上基本类型的一维或多维的数组 / 切片

## 使用示例


在任意空目录创建 `main.go` :
```go
package main

import (
    "fmt"
    "github.com/chenguofan1999/marshal"
)

type Student struct{
    Name string
    ID int
    GPA float32
    Graduated bool
    Courses []string
}

func main(){
    testStudent := Student{
        Name: "Alice",
        ID: 20180001,
        GPA: 4.2,
        Graduated: false,
        Courses: []string{
        "Data Structure",
        "Algorithm",
        "Compiler",
        "Service Computing",
        },
    }

    b, err := marshal.JsonMarshal(testStudent)
    if err != nil{
        panic(err)
    }

    fmt.Println(string(b))
}
```

此测例展示的是 `JsonMarshal` 的基本功能，将一个 struct 转换为 JSON 字节流.

**运行:**

```
$ go run main.go
{"Name":"Alice","ID":20180001,"GPA":4.2,"Graduated":false,"Courses":["Data Structure","Algorithm","Compiler","Service Computing"]}
```

### 特性：支持字段标签 `Tag` 

在结构的字段后方用以下格式可以增加字段标签 `Tag`，在转换成 JSON 字节流时会用标签内容 `TagInfo` 替代字段名。

```go
type typeName struct{
    ...
    FieldName FieldType `myTag:"TagInfo"`
    ...
}
```

**注意：** 仅识别 myTag 标签。

**示例：**

```go
package main

import (
    "fmt"
    "github.com/chenguofan1999/marshal"
)

type Student struct{
  Name string `myTag:"StudentName"`
  ID int `myTag:"ID Number"`
  GPA float32 
  Graduated bool 
  Courses []string
}

func main(){

  // TestCase 1
  testStudent := Student{
    Name: "Alice",
    ID: 20180001,
    GPA: 4.2,
    Graduated: false,
    Courses: []string{
      "Data Structure",
      "Algorithm",
      "Compiler",
      "Service Computing",
    },
  }

  b, err := marshal.JsonMarshal(testStudent)
  if err != nil{
    panic(err)
  }

  fmt.Println(string(b))

}
```

**运行：**

```
$ go run main.go
{"StudentName":"Alice","ID Number":20180001,"GPA":4.2,"Graduated":false,"Courses":["Data Structure","Algorithm","Compiler","Service Computing"]}
```

可见 "Name" 变为了 "StudentName", "ID" 变为了 "ID Number".

### 特性：支持多维数组

`JsonMarshal` 支持多维数组。

**示例：**

```go
package main

import (
    "fmt"
    "github.com/chenguofan1999/marshal"
)

type TestStruct struct {
		Ints    [][]int
		Floats  [][]float32
		Strings [][]string
	}

func main(){
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

	b, err := marshal.JsonMarshal(testCase)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```

**运行：**

```
$ go run main.go
{"Ints":[[0,1,2,3],[4,5,6,7],[8,9,10,11]],"Floats":[[0,1,2],[3,4,5],[6,7,8]],"Strings":[["aa","bb","cc","dd"],["ee","ff","gg","hh"]]}
```

此处只展示了常见的二维数组，事实上能支持任意多维的数组 / 切片。

## License

这个项目用的是 MIT License.
点击 [LICENSE](LICENSE) 查看具体内容。

## 中文 API 文档

[![GoDoc](https://img.shields.io/badge/zh--CN-REFERENCE-green?style=for-the-badge&logo=appveyor)](doc_zh_CN.md)

## 设计文档

[![GoDoc](https://img.shields.io/badge/Design-Report-critical?style=for-the-badge&logo=appveyor)](specification.md)