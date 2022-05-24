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
	desHost      []byte
	desPort      []byte
	authMode     int
	usedAuthMode uint8
}

type ServerOperation func(server *Server)

func SetServerUsername(username string) ServerOperation {
	return func(server *Server) {
		server.username = username
	}
}

func SetServerPassword(password string) ServerOperation {
	return func(server *Server) {
		server.password = password
	}
}

func SetServerAuth(authMode int) ServerOperation {
	return func(server *Server) {
		server.authMode = authMode
	}
}

func NewServer(conn net.Conn, operates ...ServerOperation) *Server {
	server := &Server{
		conn: conn,
	}

	for _, operate := range operates {
		operate(server)
	}

	return server
}

func (s *Server) HandShake() (err error) {
	err = s.handShake()
	if err != nil {
		return
	}

	if s.usedAuthMode == pass {
		err = s.authentication()
		if err != nil {
			return
		}
	}

	return s.readReqInfo()
}

func (s *Server) GetDesInfo() ([]byte, []byte) {
	return s.desHost, s.desPort
}

func (s *Server) verCheck(t byte) bool {
	if t == ver {
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
		err = errors.New(fmt.Sprintf("read len < 2, buffer: %#v", buffer))
		return
	}

	buffer = buffer[:n]
	if s.verCheck(buffer[0]) {
		err = errors.New("req version not socks5")
		return
	}

	if n < int(buffer[1])+2 {
		err = errors.New(fmt.Sprintf("buffer out of range, buffer: %#v", buffer))
		return
	}

	authModes := buffer[2 : buffer[1]+2]
	for _, authMode := range authModes {
		if authMode == noAuth && s.authMode == 1 {
			_, err = s.conn.Write(noAuthResponse)
			s.usedAuthMode = noAuth
			return
		}
		if authMode == pass {
			_, err = s.conn.Write(passAuthResponse)
			s.usedAuthMode = pass
			return
		}
	}

	_, _ = s.conn.Write(errorAuthResponse)
	err = nMETHODSError
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
		err = errors.New(fmt.Sprintf("read len < 2, buffer: %#v", buffer))
		return
	}

	usernameLen := buffer[1]
	if n < int(usernameLen)+2 {
		err = errors.Wrap(outOfRangeError, "")
		return
	}

	if s.username != string(buffer[2:2+usernameLen]) {
		err = errors.New("username is failed")
		return
	}

	passwordLen := buffer[2+usernameLen]
	if n < int(usernameLen)+2+1+int(passwordLen) {
		_, err = s.conn.Write([]byte{buffer[0], 0x01})
		err = errors.Wrap(outOfRangeError, "")
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
		err = errors.New(fmt.Sprintf("read len < 4, buffer: %#v", buffer))
		return
	}
	switch buffer[1] {
	case tcp:
		err = s.handlerTcp(buffer)
	case udp:
		err = s.handlerUdp(buffer)
	default:
		err = errors.New("unsupported transport layer protocol")
	}

	// 正常应该返回连接远程对应的ip类型+ip+port，偷懒了，直接固定参数返回了
	responseMsg := successfulFirst
	responseMsg[3] = buffer[1] // 对应协议
	if err != nil {
		responseMsg[1] = 0x00 // 失败
		_, err = s.conn.Write(responseMsg)
	} else {
		// 正常应该成功与目标ip和端口建立连接后再回复
		// 此条消息，但是为了项目结构的清晰就提前回复
		// 客服端了
		_, err = s.conn.Write(responseMsg)
	}
	return
}

func (s *Server) handlerTcp(buffer []byte) (err error) {
	n := len(buffer)
	switch buffer[3] {
	case ipv4:
		hostEndPos := 4 + net.IPv4len
		if n < hostEndPos+2 {
			err = errors.Wrap(outOfRangeError, "ipv4")
			return
		}
		s.desHost, s.desPort = buffer[4:hostEndPos], buffer[hostEndPos:hostEndPos+2]
	case domain:
		domainEndPos := 5 + buffer[4]
		if n < int(domainEndPos)+2 {
			err = errors.Wrap(outOfRangeError, "domain")
			return
		}
		s.desHost, s.desPort = buffer[5:domainEndPos], buffer[domainEndPos:domainEndPos+2]
	case ipv6:
		hostEndPos := 4 + net.IPv6len
		if n < hostEndPos+2 {
			err = errors.Wrap(outOfRangeError, "ipv6")
			return
		}
		s.desHost, s.desPort = buffer[4:hostEndPos], buffer[hostEndPos:hostEndPos+2]
	default:
		err = aTYPError
	}
	return
}

//TODO UDP流量实现
func (s *Server) handlerUdp(buffer []byte) (err error) {
	err = errors.New("UDP is not supported temporarily")
	return
}
