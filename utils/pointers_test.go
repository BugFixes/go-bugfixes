package utils_test

import (
	"testing"

	"github.com/bugfixes/go-bugfixes/utils"
	"github.com/stretchr/testify/assert"
)

func TestPointer_String(t *testing.T) {
	s := "test"
	assert.Equal(t, &s, utils.Pointer(s))
}

func TestPointer_Int(t *testing.T) {
	i := 10
	assert.Equal(t, &i, utils.Pointer(i))
}

func TestPointer_Bool(t *testing.T) {
	b := true
	assert.Equal(t, &b, utils.Pointer(b))
}

func TestPointer_Float64(t *testing.T) {
	f := 10.0
	assert.Equal(t, &f, utils.Pointer(f))
}

func TestPointer_Float32(t *testing.T) {
	f := float32(10.0)
	assert.Equal(t, &f, utils.Pointer(f))
}

func TestPointer_Int64(t *testing.T) {
	i := int64(10)
	assert.Equal(t, &i, utils.Pointer(i))
}

func TestPointer_Int32(t *testing.T) {
	i := int32(10)
	assert.Equal(t, &i, utils.Pointer(i))
}

func TestPointer_Int16(t *testing.T) {
	i := int16(10)
	assert.Equal(t, &i, utils.Pointer(i))
}

func TestPointer_Int8(t *testing.T) {
	i := int8(10)
	assert.Equal(t, &i, utils.Pointer(i))
}

func TestPointer_Uint(t *testing.T) {
	u := uint(10)
	assert.Equal(t, &u, utils.Pointer(u))
}

func TestPointer_Uint64(t *testing.T) {
	u := uint64(10)
	assert.Equal(t, &u, utils.Pointer(u))
}

func TestPointer_Uint32(t *testing.T) {
	u := uint32(10)
	assert.Equal(t, &u, utils.Pointer(u))
}

func TestPointer_Uint16(t *testing.T) {
	u := uint16(10)
	assert.Equal(t, &u, utils.Pointer(u))
}

func TestPointer_Uint8(t *testing.T) {
	u := uint8(10)
	assert.Equal(t, &u, utils.Pointer(u))
}

func TestPointer_Byte(t *testing.T) {
	b := byte(10)
	assert.Equal(t, &b, utils.Pointer(b))
}

func TestPointer_Rune(t *testing.T) {
	r := rune(10)
	assert.Equal(t, &r, utils.Pointer(r))
}

func TestPointer_Complex64(t *testing.T) {
	c := complex64(10)
	assert.Equal(t, &c, utils.Pointer(c))
}

func TestPointer_Complex128(t *testing.T) {
	c := complex128(10)
	assert.Equal(t, &c, utils.Pointer(c))
}
