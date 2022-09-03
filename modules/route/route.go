package route

import (
	pubConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/route/direct"
	"agile-proxy/modules/route/dynamic"
	"agile-proxy/modules/route/ha"
	sysJson "encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Route interface {
	Run() (err error)
	Transport(conn net.Conn, host, port []byte) (err error)
	HttpTransport(w http.ResponseWriter, r *http.Request) (err error)
	Close() (err error)
}

func Factory(configs []sysJson.RawMessage) {
	for _, config := range configs {
		var err error
		var transport Route
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
			route[transportName] = transport
		}
	}
}
