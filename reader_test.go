package bin

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Point struct {
	X int
	Y int
}

func TestReader_ReadVar(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriter(buf)
	w.WriteVar(uint64(123))
	w.WriteVar("abc")
	w.WriteVar(3.1415)
	w.WriteVar([]byte{5, 6, 7})
	w.WriteVar(Point{88, 99})

	var (
		i int
		s string
		f float64
		b []byte
		p Point
	)
	r := NewReader(buf)
	r.ReadVar(&i)
	r.ReadVar(&s)
	r.ReadVar(&f)
	r.ReadVar(&b)
	r.ReadVar(&p)

	assert.Equal(t, 123, i)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 3.1415, f)
	assert.Equal(t, []byte{5, 6, 7}, b)
	assert.Equal(t, Point{88, 99}, p)
}

func TestReader_ReadVarInt(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVarInt(0x1234)

	r := w.Reader
	iDec, err := r.ReadVarInt()

	assert.NoError(t, err)
	assert.Equal(t, 0x1234, iDec)
}

func TestReader_ReadFloat64(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(-0.123456789)

	r := w.Reader
	res, err := r.ReadFloat64()

	assert.NoError(t, err)
	assert.Equal(t, -0.123456789, res)
}

func TestReader_ReadFloat32(t *testing.T) {
	w := NewBuffer(nil)
	w.WriteVar(float32(-1. / 3))

	r := w.Reader
	res, err := r.ReadFloat32()

	assert.NoError(t, err)
	assert.Equal(t, float32(-1./3), res)
}

func TestReader_ReadTime(t *testing.T) {
	time.Local = time.UTC
	w := NewBuffer(nil)
	w.WriteVar(time.Date(2016, 07, 06, 18, 24, 45, 0, time.Local))

	r := w.Reader
	res, err := r.ReadTime()

	assert.NoError(t, err)
	assert.Equal(t, "2016-07-06 18:24:45 UTC", res.Format("2006-01-02 15:04:05 MST"))
}

func TestReader_ReadTime32(t *testing.T) {
	time.Local = time.UTC
	w := NewBuffer(nil)
	w.WriteTime32(time.Date(2016, 07, 06, 18, 24, 45, 0, time.Local))

	r := w.Reader
	res, err := r.ReadTime32()

	assert.NoError(t, err)
	assert.Equal(t, "2016-07-06 18:24:45 UTC", res.Format("2006-01-02 15:04:05 MST"))
}

func TestReader_ReadIntSlice(t *testing.T) {
	w := NewBuffer(nil)
	for i := 0; i < 100; i++ {
		w.WriteVar(i)
	}
	r := w.Reader

	for i := 0; i < 100; i++ {
		var j int
		r.ReadVar(&j)
		assert.Equal(t, j, i)
	}
}

func TestReader_SetReadLimit(t *testing.T) {
	arr := make([]byte, 100)
	w := NewBuffer(nil)
	w.WriteVar(arr)

	r := w.Reader
	r.SetReadLimit(101) // 100 + 1 byte of length
	res, err := r.ReadBytes()

	assert.NoError(t, err)
	assert.Equal(t, 100, len(res))
}

func TestReader_ReadLimit_Fail(t *testing.T) {
	arr := make([]byte, 100)
	w := NewBuffer(nil)
	w.WriteVar(arr)

	r := w.Reader
	r.SetReadLimit(99)
	_, err := r.ReadBytes()

	assert.Error(t, err)
}
