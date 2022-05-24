package dialer

import (
	"github.com/pkg/errors"
	"net"
	pConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/dialer/direct"
	"strings"
	"time"
)

type Dialer interface {
	Dial(network string, host, port string) (conn net.Conn, err error)
	DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)
}

type BaseDialer struct {
	Name  string
	Type  string
	IFace string
}

func Factory(configs []string) {
	for _, config := range configs {
		var err error
		var dialer Dialer
		switch strings.ToLower(json.Get([]byte(config), "type").ToString()) {
		case pConfig.Socks5:
		case pConfig.Ssh:
		case pConfig.Direct:
			dialer, err = direct.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("server init failed: %v", err)
			continue
		}

		dialerName := json.Get([]byte(config), "name").ToString()
		if dialerName != "" {
			dialers[dialerName] = dialer
		}
	}
}
