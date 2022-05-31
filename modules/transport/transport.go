package transport

import (
	pubConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/transport/direct"
	official "encoding/json"
	"errors"
	"net"
	"strings"
)

type Transport interface {
	Transport(conn net.Conn, host, port []byte) (err error)
	Close() (err error)
}

func Factory(configs []official.RawMessage) {
	for _, config := range configs {
		var err error
		var transport Transport
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case pubConfig.Direct:
			transport, err = direct.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("%#v", err)
			continue
		}

		transportName := json.Get(config, "name").ToString()
		if transportName != "" {
			transports[transportName] = transport
		}
	}
}
