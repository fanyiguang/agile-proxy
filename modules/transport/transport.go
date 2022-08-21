package transport

import (
	pubConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/transport/direct"
	"agile-proxy/modules/transport/dynamic"
	"agile-proxy/modules/transport/ha"
	sysJson "encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Transport interface {
	Run() (err error)
	Transport(conn net.Conn, host, port []byte) (err error)
	Close() (err error)
}

func Factory(configs []sysJson.RawMessage) {
	for _, config := range configs {
		var err error
		var transport Transport
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case pubConfig.Direct:
			transport, err = direct.New(config)
		case pubConfig.Dynamic:
			transport, err = dynamic.New(config)
		case pubConfig.Ha:
			transport, err = ha.New(config)
		default:
			err = errors.New(fmt.Sprintf("type is invalid %v", json.Get(config, "type").ToString()))
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
