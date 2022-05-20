package socks5

import "github.com/pkg/errors"

const (
	ver uint8 = 0x05

	noAuth uint8 = 0x00
	GSSAPI uint8 = 0x01
	pass   uint8 = 0x02

	tcp uint8 = 0x01
	udp uint8 = 0x03

	ipv4   uint8 = 0x01
	domain uint8 = 0x03
	ipv6   uint8 = 0x04
)

var (
	noAuthResponse    = []byte{0x05, 0x00}
	passAuthResponse  = []byte{0x05, 0x02}
	errorAuthResponse = []byte{0x05, 0xff}

	successfulFirst = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} //链接成功报文

	nMETHODSError   = errors.New("wrong NMETHODS")
	aTYPError       = errors.New("wrong ATYP")
	outOfRangeError = errors.New("slice out of range")
)
