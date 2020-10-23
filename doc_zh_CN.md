
## func `JsonMarshal`

```go
func JsonMarshal(v interface{}) ([]byte, error)
```
JsonMarshal 接受一个任意的 struct，返回其 JSON 格式的字节流。
是本包用于**序列化**的主要函数。


## func `parse`

```go
func Parse(valField reflect.Value) string
```

Parse 接受（来自一个结构的）一个域的 reflect.Value，返回其 JSON 字符串形式。是本包中服务于 `JsonMarshal` 的辅助函数。