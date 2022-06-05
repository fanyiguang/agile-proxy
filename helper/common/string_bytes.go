package common

// https://github.com/qcrao/Go-Questions/issues/7

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func GetBytesLen(b []byte) int {
	return (*reflect.SliceHeader)(unsafe.Pointer(&b)).Len
}

func IntToBytes(n int, b byte) ([]byte, error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 3, 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}

func CopyBytes(des []byte) []byte {
	clone := make([]byte, len(des))
	copy(clone, des)
	des = clone
	return des
}
