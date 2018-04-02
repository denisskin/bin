package bin

type Encoder interface {
	Encode() []byte
}

type Decoder interface {
	Decode([]byte) error
}

func Encode(vv ...interface{}) []byte {
	w := NewBuffer(nil)
	for _, v := range vv {
		w.WriteVar(v)
	}
	return w.Bytes()
}

func Decode(data []byte, vv ...interface{}) error {
	r := NewBuffer(data)
	for _, v := range vv {
		r.ReadVar(v)
	}
	return r.Error()
}

type binWriter interface {
	BinWrite(writer *Writer)
}

type binReader interface {
	BinRead(reader *Reader)
}
