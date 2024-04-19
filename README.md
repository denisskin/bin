# bin
Go library. Simple reader and writer for marshaling binary data.

```go
//---- write
buf := bytes.NewBuffer(nil)
w := NewWriter(buf)
w.WriteVar(uint64(123))
w.WriteVar("abc")
w.WriteVar(3.1415)

//----- read
var (
    i  int
    s  string
    f  float64
)
r := NewReader(buf)
r.ReadVar(&i) // 123
r.ReadVar(&s) // "abc"
r.ReadVar(&f) // 3.1415
```

Encode, decode var int
```go
w := NewBuffer(nil)
w.WriteVarInt(0x1234)

r := w.Reader
iDec, err := r.ReadVarInt() // -> 0x1234, nil
```