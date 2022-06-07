package socks5

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

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

func ChangeStrPort(bPort []byte) string {
	// 两个字节刚好够装所有端口
	if len(bPort) < 2 {
		return ""
	}
	return strconv.Itoa(int(bPort[0])<<8 | int(bPort[1]))
}
