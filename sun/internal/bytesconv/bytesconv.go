package bytesconv

import "unsafe"

// StringToBytes 使用这个工具和直接Write([]bytes())相比，会少一些拷贝，性能更高
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
