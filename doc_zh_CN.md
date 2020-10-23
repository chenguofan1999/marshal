
## func `JsonMarshal`

```go
func JsonMarshal(v interface{}) ([]byte, error)
```
JsonMarshal 接受一个任意的 struct，返回其 JSON 格式的字节流。
是本包用于**序列化**的主要函数。

