package dialer

import (
	pConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/dialer/direct"
	official "encoding/json"
	"github.com/pkg/errors"
	"net"
	"strings"
	"time"
)

type Dialer interface {
	Dial(network string, host, port string) (conn net.Conn, err error)
	DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)
}

func Factory(configs []official.RawMessage) {
	for _, config := range configs {
		var err error
		var dialer Dialer
		switch strings.ToLower(json.Get(config, "type").ToString()) {
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
