package transport

import (
	"errors"
	"io"
	"net"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/transport/direct"
	"strings"
	"time"
)

type Transport interface {
	Transport(conn net.Conn, host, port []byte) (err error)
	Close() (err error)
}

type DnsInfo struct {
	Server   string `json:"server"`
	LocalDns bool   `json:"local_dns"`
}

type BaseTransport struct {
	Type       string
	Name       string
	ClientName string
}

func (t BaseTransport) MyCopy(sConn net.Conn, cConn net.Conn) {
	cDoneCh, sDoneCh := make(chan struct{}), make(chan struct{})
	go t.myCopyN(sConn, cConn, cDoneCh)
	go t.myCopyN(cConn, sConn, sDoneCh)
	t.wait(cDoneCh, sDoneCh)
}

func (t BaseTransport) myCopyN(det net.Conn, src net.Conn, done chan struct{}) {
	for {
		_ = src.SetReadDeadline(time.Now().Add(2 * time.Minute))
		_ = det.SetWriteDeadline(time.Now().Add(2 * time.Minute))
		_, err := io.CopyN(det, src, 1024*32)
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
	}
	close(done)
}

func (t BaseTransport) wait(cDoneCh, sDoneCh chan struct{}) {
	<-cDoneCh
	<-sDoneCh
}

func Factory(configs []string) {
	for _, config := range configs {
		var err error
		var transport Transport
		switch strings.ToLower(json.Get([]byte(config), "type").ToString()) {
		case Normal:
			transport, err = direct.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("%#v", err)
			continue
		}

		transportName := json.Get([]byte(config), "name").ToString()
		if transportName != "" {
			transports[transportName] = transport
		}
	}
}
