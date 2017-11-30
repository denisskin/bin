package bin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuffer(t *testing.T) {

	w := NewWriter(bytes.NewBuffer(nil))

	assert.True(t, w != nil)
}

func TestWriter_WriteVar(t *testing.T) {
	w := NewBuffer(nil)

	w.WriteVar(0)
	w.WriteVar(13)
	w.WriteVar(255)
	w.WriteVar(256)
	w.WriteVar(-13)
	w.WriteVar(0x01020304050607)
	var max64 uint64 = 0xffffffffffffffff
	w.WriteVar(max64)
	w.WriteVar(0.3)

	assert.Equal(t, []byte{
		0,
		13,
		0x81, 0xff,
		0x82, 1, 0,
		0xc1, 13,
		0x87, 1, 2, 3, 4, 5, 6, 7,
		0xc1, 1,
		0x3f, 0xd3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
	}, w.Bytes())
}

func TestWriter_WriteString(t *testing.T) {
	w := NewBuffer(nil)

	w.WriteVar("")
	w.WriteVar("Abc")

	assert.Equal(t, []byte{0, 3, 'A', 'b', 'c'}, w.Bytes())
}
