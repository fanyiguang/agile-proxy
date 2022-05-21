package common

// https://github.com/qcrao/Go-Questions/issues/7

import (
	"reflect"
	"unsafe"
)

func BytesToStr(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}

	return *(*string)(unsafe.Pointer(&sh))
}

func StrToBytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
