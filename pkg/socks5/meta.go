package socks5

import "github.com/pkg/errors"

const (
	VER uint8 = 0x05

	NoAuth uint8 = 0x00
	GSSAPI uint8 = 0x01
	Pass   uint8 = 0x02

	tcp uint8 = 0x01
	udp uint8 = 0x03
)

var (
	noAuthResponse    = []byte{0x05, 0x00}
	passAuthResponse  = []byte{0x05, 0x02}
	errorAuthResponse = []byte{0x05, 0xff}

	successfulFirst  = []byte{0x05, 0x00, 0x00}                   //链接成功报文前半段
	successfulSecond = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00} //链接成功报文后半段

	authProtocolError = errors.New("Authentication protocol not supported")
)
