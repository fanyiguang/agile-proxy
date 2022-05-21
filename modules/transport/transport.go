package transport

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"net"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/transport/direct"
	"strings"
)

var (
	transports = make(map[string]Transport)
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Transport interface {
	Transport(conn net.Conn, host, port []byte) (err error)
	Close() (err error)
}

func GetTransport(name string) (t Transport) {
	return transports[name]
}

func GetAllTransports() map[string]Transport {
	return transports
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
