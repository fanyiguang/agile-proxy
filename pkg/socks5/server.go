package socks5

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
)

type Server struct {
	conn         net.Conn
	username     string
	password     string
	authMode     int
	usedAuthMode uint8
}

type Operation func(server *Server)

func SetUsername(username string) Operation {
	return func(server *Server) {
		server.username = username
	}
}

func SetPassword(password string) Operation {
	return func(server *Server) {
		server.password = password
	}
}

func SetAuth(authMode int) Operation {
	return func(server *Server) {
		server.authMode = authMode
	}
}

func NewServer(conn net.Conn, operates ...Operation) *Server {
	server := &Server{
		conn: conn,
	}

	for _, operate := range operates {
		operate(server)
	}

	return server
}

func (s *Server) Run() (err error) {
	err = s.handShake()
	if err != nil {
		return
	}

	if s.usedAuthMode == Pass {
		err = s.authentication()
		if err != nil {
			return
		}
	}

	s.readReqInfo()
	return
}

func (s *Server) verCheck(t byte) bool {
	if t == VER {
		return true
	} else {
		return false
	}
}

func (s *Server) handShake() (err error) {
	buffer := make([]byte, 16)
	n, err := s.conn.Read(buffer)
	if err != nil {
		err = errors.Wrap(err, "reader.Read")
		return
	}

	if n < 2 {
		err = errors.New(fmt.Sprintf("read len < 2, buffer: %x", buffer))
		return
	}

	buffer = buffer[:n]
	if s.verCheck(buffer[0]) {
		err = errors.New("req version not socks5")
		return
	}

	if n < int(buffer[1])+2 {
		err = errors.New(fmt.Sprintf("buffer out of range, buffer: %x", buffer))
		return
	}

	authModes := buffer[2 : buffer[1]+2]
	for _, authMode := range authModes {
		if authMode == NoAuth && s.authMode == 1 {
			_, err = s.conn.Write(noAuthResponse)
			s.usedAuthMode = NoAuth
			return
		}
		if authMode == Pass {
			_, err = s.conn.Write(passAuthResponse)
			s.usedAuthMode = Pass
			return
		}
	}

	_, _ = s.conn.Write(errorAuthResponse)
	err = authProtocolError
	return
}

func (s *Server) authentication() (err error) {
	buffer := make([]byte, 512)
	n, err := s.conn.Read(buffer)
	if err != nil {
		err = errors.Wrap(err, "reader.Read")
		return
	}

	if n < 2 {
		err = errors.New(fmt.Sprintf("read len < 2, buffer: %x", buffer))
		return
	}

	usernameLen := buffer[1]
	if n < int(usernameLen)+1 {
		err = errors.Wrap(errors.New("slice out of range"), "")
		return
	}

	if s.username != string(buffer[2:2+usernameLen]) {
		err = errors.New("username is failed")
		return
	}

	passwordLen := buffer[2+usernameLen]
	if n < int(usernameLen)+1+1+int(passwordLen) {
		_, err = s.conn.Write([]byte{buffer[0], 0x01})
		err = errors.Wrap(errors.New("slice out of range"), "")
		return
	}

	if s.password != string(buffer[2+usernameLen+1:2+usernameLen+1+passwordLen]) {
		_, err = s.conn.Write([]byte{buffer[0], 0x02})
		err = errors.New("password is failed")
		return
	}

	_, err = s.conn.Write([]byte{buffer[0], 0x00})
	return
}

func (s *Server) readReqInfo() (err error) {
	buffer := make([]byte, 512)
	n, err := s.conn.Read(buffer)
	if err != nil {
		err = errors.Wrap(err, "reader.Read")
		return
	}

	if n < 4 {
		err = errors.New(fmt.Sprintf("read len < 4, buffer: %x", buffer))
		return
	}
	switch buffer[1] {
	case tcp:

		//s.conn.Write([]byte{successfulFirst, 0x00})
	case udp:

	default:
		err = errors.New("unsupported transport layer protocol")
	}
	return
}
