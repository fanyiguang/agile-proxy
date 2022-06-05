package dialer

import (
	pConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/dialer/direct"
	"agile-proxy/modules/dialer/http"
	"agile-proxy/modules/dialer/https"
	"agile-proxy/modules/dialer/socks5"
	"agile-proxy/modules/dialer/ssh"
	sysJson "encoding/json"
	"github.com/pkg/errors"
	"net"
	"strings"
	"time"
)

type Dialer interface {
	Dial(network string, host, port string) (conn net.Conn, err error)
	DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)
}

func Factory(configs []sysJson.RawMessage) {
	for _, config := range configs {
		var err error
		var dialer Dialer
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case pConfig.Socks5:
			dialer, err = socks5.New(config)
		case pConfig.Ssh:
			dialer, err = ssh.New(config)
		case pConfig.Https:
			dialer, err = https.New(config)
		case pConfig.Http:
			dialer, err = http.New(config)
		case pConfig.Direct:
			dialer, err = direct.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("server init failed: %v", err)
			continue
		}

		dialerName := json.Get(config, "name").ToString()
		if dialerName != "" {
			dialers[dialerName] = dialer
		}
	}
}
