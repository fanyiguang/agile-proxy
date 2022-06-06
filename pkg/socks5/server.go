package socks5

import (
	"agile-proxy/helper/common"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
)

type Server struct {
	bufferPool sync.Pool
	username   string
	password   string
	authMode   int
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

func NewServer(operates ...ServerOperation) *Server {
	server := &Server{
		bufferPool: sync.Pool{New: func() any {
			return make([]byte, 512)
		}},
	}
	for _, operate := range operates {
		operate(server)
	}
	return server
}

func (s *Server) HandShake(conn net.Conn) (desHost, desPort []byte, err error) {
	usedAuthMode, err := s.handShake(conn)
	if err != nil {
		return
	}

	if usedAuthMode == pass {
		err = s.authentication(conn)
		if err != nil {
			return
		}
	}

	return s.readReqInfo(conn)
}

func (s *Server) verCheck(t byte) bool {
	if t == ver {
		return true
	} else {
		return false
	}
}

func (s *Server) handShake(conn net.Conn) (usedAuthMode uint8, err error) {
	// 理论最大值为1+1+255，但是一般不会到那么大
	// 为了节省空间设置为8，有特殊情况再做修改
	buffer := make([]byte, 8)
	//buffer := s.bufferPool.Get().([]byte)
	//defer s.bufferPool.Put(buffer)
	n, err := conn.Read(buffer)
	if err != nil {
		err = errors.Wrap(err, "reader.Read")
		return
	}

	if n < 2 {
		err = errors.New(fmt.Sprintf("read len < 2, buffer: %#v", buffer))
		return
	}

	buffer = buffer[:n]
	if !s.verCheck(buffer[0]) {
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
			_, err = conn.Write(noAuthResponse)
			usedAuthMode = noAuth
			return
		}
		if authMode == pass {
			_, err = conn.Write(passAuthResponse)
			usedAuthMode = pass
			return
		}
	}

	_, _ = conn.Write(errorAuthResponse)
	err = nMETHODSError
	return
}

func (s *Server) authentication(conn net.Conn) (err error) {
	buffer := s.bufferPool.Get().([]byte)
	defer s.bufferPool.Put(buffer)
	n, err := conn.Read(buffer)
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
		_, err = conn.Write([]byte{buffer[0], 0x01})
		err = errors.Wrap(outOfRangeError, "")
		return
	}

	if s.password != string(buffer[2+usernameLen+1:2+usernameLen+1+passwordLen]) {
		_, err = conn.Write([]byte{buffer[0], 0x02})
		err = errors.New("password is failed")
		return
	}

	_, err = conn.Write([]byte{buffer[0], 0x00})
	return
}

func (s *Server) readReqInfo(conn net.Conn) (desHost, desPort []byte, err error) {
	buffer := s.bufferPool.Get().([]byte)
	defer s.bufferPool.Put(buffer)
	n, err := conn.Read(buffer)
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
		desHost, desPort, err = s.handlerTcp(buffer, n)
	case udp:
		desHost, desPort, err = s.handlerUdp(buffer, n)
	default:
		err = errors.New("unsupported transport layer protocol")
	}

	// 因为使用了buffer复用并且传递的是切片，会有脏数据的问题。
	// 所以需要复制一份 host and port，彻底与原来的底层数组
	// 脱离关系。
	desHost, desPort = common.CopyBytes(desHost), common.CopyBytes(desPort)

	// 正常应该返回连接远程对应的ip类型+ip+port，偷懒了，直接固定参数返回了
	responseMsg := successfulFirst
	responseMsg[3] = buffer[1] // 对应协议
	if err != nil {
		responseMsg[1] = 0x00 // 失败
		_, err = conn.Write(responseMsg)
	} else {
		// 正常应该成功与目标ip和端口建立连接后再回复
		// 此条消息，但是为了项目结构的清晰就提前回复
		// 客户端了
		_, err = conn.Write(responseMsg)
	}
	return
}

func (s *Server) handlerTcp(buffer []byte, n int) (desHost, desPort []byte, err error) {
	switch buffer[3] {
	case ipv4:
		hostEndPos := 4 + net.IPv4len
		if n < hostEndPos+2 {
			err = errors.Wrap(outOfRangeError, "ipv4")
			return
		}
		desHost, desPort = buffer[4:hostEndPos], buffer[hostEndPos:hostEndPos+2]
	case domain:
		domainEndPos := 5 + buffer[4]
		if n < int(domainEndPos)+2 {
			err = errors.Wrap(outOfRangeError, "domain")
			return
		}
		desHost, desPort = buffer[5:domainEndPos], buffer[domainEndPos:domainEndPos+2]
	case ipv6:
		hostEndPos := 4 + net.IPv6len
		if n < hostEndPos+2 {
			err = errors.Wrap(outOfRangeError, "ipv6")
			return
		}
		desHost, desPort = buffer[4:hostEndPos], buffer[hostEndPos:hostEndPos+2]
	default:
		err = aTYPError
	}
	return
}

//TODO UDP流量实现
func (s *Server) handlerUdp(buffer []byte, n int) (desHost, desPort []byte, err error) {
	err = errors.New("UDP is not supported temporarily")
	return
}
