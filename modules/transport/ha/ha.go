package ha

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport/base"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
	"time"
)

type ha struct {
	baseTransport base.Transport
	clients       []client.Client
}

func (h *ha) Transport(cConn net.Conn, host, port []byte) (err error) {
	if h.clients != nil {
		host, err = h.baseTransport.Dns.GetHost(host)
		if err != nil {
			return
		}

		sConn, err := h.getConn(host, port)
		if err != nil {
			return err
		}

		defer sConn.Close()
		h.baseTransport.AsyncSendMsgToIpc(fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		h.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("client is nil")
	}
	return
}

func (h *ha) Close() (err error) {
	return
}

func (h *ha) getConn(host, port []byte) (conn net.Conn, err error) {
	connCh := make(chan net.Conn)
	for _, c := range h.clients {
		_client := c
		Go.Go(func() {
			var sConn net.Conn
			sConn, err = _client.Dial("tcp", host, port)
			if err != nil {
				log.DebugF("ha client dial failed: %v %v %v %v", err, host, port, _client.Name())
				return
			}

			select {
			case connCh <- sConn:
			default:
				_ = sConn.Close()
			}
		})
	}

	select {
	case conn = <-connCh:
	case <-time.After(time.Second * 15):
		err = errors.New("get conn timeout")
	}
	return
}

func New(jsonConfig json.RawMessage) (obj *ha, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		marshalJSON, _ := jsonConfig.MarshalJSON()
		err = errors.Wrap(err, common.BytesToStr(marshalJSON))
		return
	}

	if !strings.Contains(config.DnsInfo.Server, ":") {
		config.DnsInfo.Server = net.JoinHostPort(config.DnsInfo.Server, "53")
	}

	obj = &ha{
		baseTransport: base.Transport{
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Dns: plugin.Dns{
				Server:   config.DnsInfo.Server,
				LocalDns: config.DnsInfo.LocalDns,
			},
			BufferPool: sync.Pool{
				New: func() any {
					return make([]byte, 1024*32)
				},
			},
		},
	}

	if config.ClientNames != "" {
		clientNames := strings.Split(config.ClientNames, ",")
		for _, clientName := range clientNames {
			_client := client.GetClient(clientName)
			if _client != nil {
				obj.clients = append(obj.clients, _client)
			}
		}
	}

	return
}
