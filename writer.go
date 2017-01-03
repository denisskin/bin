package bin

import (
	"bytes"
	"io"
	"math"
	"time"
)

type Writer struct {
	wr         io.Writer
	err        error
	CntWritten uint64
}

func NewWriter(w io.Writer) *Writer {
	if w == nil {
		w = bytes.NewBuffer(nil)
	}
	return &Writer{wr: w}
}

func (w *Writer) Error() error {
	return w.err
}

func (w *Writer) SetError(err error) {
	if err != nil {
		w.err = err
	}
}

func (w *Writer) write(bb []byte) error {
	w.Write(bb)
	return w.err
}

func (w *Writer) Write(bb []byte) (n int, err error) {
	if buf, ok := w.wr.(*bytes.Buffer); ok {
		n, err = buf.Write(bb)
	} else {
		var n64 int64
		n64, err = io.Copy(w.wr, bytes.NewBuffer(bb))
		n = int(n64)
	}
	w.CntWritten += uint64(n)
	w.SetError(err)
	return
}

//----------- fixed types --------------
func (w *Writer) WriteByte(b byte) error {
	return w.write([]byte{b})
}

func (w *Writer) WriteUint8(i uint8) error {
	return w.WriteByte(byte(i))
}

func (w *Writer) WriteUint16(i uint16) error {
	return w.write(Uint16ToBytes(i))
}

func (w *Writer) WriteUint32(i uint32) error {
	return w.write(Uint32ToBytes(i))
}

func (w *Writer) WriteUint64(i uint64) error {
	return w.write(Uint64ToBytes(i))
}

func (w *Writer) WriteFloat32(f float32) error {
	return w.write(Uint32ToBytes(math.Float32bits(f)))
}

func (w *Writer) WriteFloat64(f float64) error {
	return w.write(Uint64ToBytes(math.Float64bits(f)))
}

func (w *Writer) WriteTime(t time.Time) error {
	return w.write(Uint64ToBytes(uint64(t.UnixNano())))
}

func (w *Writer) WriteTime32(t time.Time) error {
	return w.write(Uint32ToBytes(uint32(t.Unix())))
}

func (w *Writer) WriteBool(f bool) error {
	if f {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

//----------- var types ----------------
func (w *Writer) WriteVarInt(num int) error {
	return w.writeVarInt(int64(num))
}

func (w *Writer) writeVarInt(i int64) error {
	if i >= 0 && i < 128 {
		return w.write([]byte{byte(i)})
	}
	var h byte = 0x80
	if i < 0 {
		h |= 0x40
		i = -i
	}
	buf := make([]byte, 8)
	var n byte = 0
	for i > 0 {
		n++
		buf[8-n] = byte(i)
		i >>= 8
	}
	return w.write(append([]byte{h | n}, buf[8-n:]...))
}

func (w *Writer) WriteSliceBytes(bb [][]byte) error {
	w.WriteVarInt(len(bb))
	for _, d := range bb {
		w.WriteBytes(d)
	}
	return w.err
}

func (w *Writer) WriteBytes(bb []byte) error {
	w.WriteVarInt(len(bb))
	w.Write(bb)
	return w.err
}

func (w *Writer) WriteString(s string) error {
	w.WriteBytes([]byte(s))
	return w.err
}

func (w *Writer) WriteSliceString(ss []string) error {
	w.WriteVarInt(len(ss))
	for _, s := range ss {
		if w.WriteString(s) != nil {
			break
		}
	}
	return w.err
}

func (w *Writer) WriteError(err error) error {
	return w.WriteString(err.Error())
}

func (w *Writer) WriteVar(val interface{}) error {
	switch v := val.(type) {
	case nil:
		w.Write([]byte{0})

	case int:
		w.writeVarInt(int64(v))
	case int8:
		w.writeVarInt(int64(v))
	case int16:
		w.writeVarInt(int64(v))
	case int32:
		w.writeVarInt(int64(v))
	case int64:
		w.writeVarInt(int64(v))

	case uint:
		w.writeVarInt(int64(v))
	case uint8:
		w.writeVarInt(int64(v))
	case uint16:
		w.writeVarInt(int64(v))
	case uint32:
		w.writeVarInt(int64(v))
	case uint64:
		w.writeVarInt(int64(v))

	case float32:
		w.WriteFloat32(v)
	case float64:
		w.WriteFloat64(v)
	case time.Time:
		w.WriteTime(v)
	case bool:
		w.WriteBool(v)

	case string:
		w.WriteString(v)
	case []string:
		w.WriteSliceString(v)
	case []byte:
		w.WriteBytes(v)
	case [][]byte:
		w.WriteSliceBytes(v)

	case BinEncoder:
		if data, err := Encode(v); err == nil {
			w.WriteBytes(data)
		} else {
			w.SetError(err)
		}

	case error:
		w.WriteError(v)

	default:
		panic("Unkonown type")
	}
	return w.err
}
