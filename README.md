# bin
Go library. Simple reader and writer for marshaling binary data.

```go
//---- write
buf := bytes.NewBuffer(nil)
w := NewWriter(buf)
w.WriteVar(uint64(123))
w.WriteVar("abc")
w.WriteVar(3.1415)
w.WriteVar([]byte{5, 6, 7})
w.WriteVar([]int{7, 8, 9})

//----- read
var (
    i  int
    s  string
    f  float64
    bb []byte
    ii []int
)
r := NewReader(buf)
r.ReadVar(&i)   // 123
r.ReadVar(&s)   // "abc"
r.ReadVar(&f)   // 3.1415
r.ReadVar(&bb)  // []byte{5, 6, 7}
r.ReadVar(&ii)  // []int{7, 8, 9}
```

Encode, decode var int
```go
w := NewBuffer(nil)
w.WriteVarInt(0x1234)

r := w.Reader
iDec, err := r.ReadVarInt() // -> 0x1234, nil
```