package utils

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}

func BoolPointer(b bool) *bool {
	return &b
}

func Float64Pointer(f float64) *float64 {
	return &f
}

func Float32Pointer(f float32) *float32 {
	return &f
}

func Int64Pointer(i int64) *int64 {
	return &i
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func Int16Pointer(i int16) *int16 {
	return &i
}

func Int8Pointer(i int8) *int8 {
	return &i
}

func UintPointer(u uint) *uint {
	return &u
}

func Uint64Pointer(u uint64) *uint64 {
	return &u
}

func Uint32Pointer(u uint32) *uint32 {
	return &u
}

func Uint16Pointer(u uint16) *uint16 {
	return &u
}

func Uint8Pointer(u uint8) *uint8 {
	return &u
}

func BytePointer(b byte) *byte {
	return &b
}

func RunePointer(r rune) *rune {
	return &r
}

func Complex64Pointer(c complex64) *complex64 {
	return &c
}

func Complex128Pointer(c complex128) *complex128 {
	return &c
}

func Pointer[T any](v T) *T {
	return &v
}
