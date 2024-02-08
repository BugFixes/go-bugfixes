package utils_test

import (
	"github.com/bugfixes/go-bugfixes/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringPointer(t *testing.T) {
	s := "test"
	sp := utils.StringPointer(s)
	assert.Equal(t, &s, sp)
}

func TestIntPointer(t *testing.T) {
	i := 10
	ip := utils.IntPointer(i)
	assert.Equal(t, &i, ip)
}

func TestBoolPointer(t *testing.T) {
	b := true
	bp := utils.BoolPointer(true)
	assert.Equal(t, &b, bp)
}

func TestFloat64Pointer(t *testing.T) {
	f := 10.0
	fp := utils.Float64Pointer(f)
	assert.Equal(t, &f, fp)
}

func TestFloat32Pointer(t *testing.T) {
	f := float32(10.0)
	fp := utils.Float32Pointer(f)
	assert.Equal(t, &f, fp)
}

func TestInt64Pointer(t *testing.T) {
	i := int64(10)
	ip := utils.Int64Pointer(i)
	assert.Equal(t, &i, ip)
}

func TestInt32Pointer(t *testing.T) {
	i := int32(10)
	ip := utils.Int32Pointer(i)
	assert.Equal(t, &i, ip)
}

func TestInt16Pointer(t *testing.T) {
	i := int16(10)
	ip := utils.Int16Pointer(i)
	assert.Equal(t, &i, ip)
}

func TestInt8Pointer(t *testing.T) {
	i := int8(10)
	ip := utils.Int8Pointer(i)
	assert.Equal(t, &i, ip)
}

func TestUintPointer(t *testing.T) {
	u := uint(10)
	up := utils.UintPointer(u)
	assert.Equal(t, &u, up)
}

func TestUint64Pointer(t *testing.T) {
	u := uint64(10)
	up := utils.Uint64Pointer(u)
	assert.Equal(t, &u, up)
}

func TestUint32Pointer(t *testing.T) {
	u := uint32(10)
	up := utils.Uint32Pointer(u)
	assert.Equal(t, &u, up)
}

func TestUint16Pointer(t *testing.T) {
	u := uint16(10)
	up := utils.Uint16Pointer(u)
	assert.Equal(t, &u, up)
}

func TestUint8Pointer(t *testing.T) {
	u := uint8(10)
	up := utils.Uint8Pointer(u)
	assert.Equal(t, &u, up)
}

func TestBytePointer(t *testing.T) {
	b := byte(10)
	bp := utils.BytePointer(b)
	assert.Equal(t, &b, bp)
}

func TestRunePointer(t *testing.T) {
	r := rune(10)
	rp := utils.RunePointer(r)
	assert.Equal(t, &r, rp)
}

func TestComplex64Pointer(t *testing.T) {
	c := complex64(10)
	cp := utils.Complex64Pointer(c)
	assert.Equal(t, &c, cp)
}

func TestComplex128Pointer(t *testing.T) {
	c := complex128(10)
	cp := utils.Complex128Pointer(c)
	assert.Equal(t, &c, cp)
}

func TestPointer(t *testing.T) {
	i := 10
	ip := utils.Pointer(i)
	assert.Equal(t, &i, ip)
}
