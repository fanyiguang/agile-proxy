package dynamic

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport/base"
	"agile-proxy/modules/transport/dynamic/rule"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
)

type dynamic struct {
	baseTransport base.Transport
	clients       []client.Client // 动态类型的传输器客户端可以为多个
	rule          rule.Rule
	clientLen     int
}

func (d *dynamic) Transport(cConn net.Conn, host, port []byte) (err error) {
	if d.clients != nil {
		host, err = d.baseTransport.GetHost(host)
		if err != nil {
			return
		}

		var sConn net.Conn
		sConn, err = d.clients[d.getClientIndex()].Dial("tcp", host, port)
		if err != nil {
			return
		}

		defer sConn.Close()
		d.baseTransport.AsyncSendMsgToIpc(fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		d.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("Client is nil")
	}
	return
}

func (d *dynamic) Close() (err error) {
	for key, _client := range d.clients {
		if _client != nil {
			err = _client.Close()
			log.DebugF("dynamic close failed: %v %v", err, key)
		}
	}
	return
}

func (d *dynamic) getClientIndex() (idx int) {
	return d.rule.Intn(d.clientLen)
}

func New(jsonConfig json.RawMessage) (obj *dynamic, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		marshalJSON, _ := jsonConfig.MarshalJSON()
		err = errors.Wrap(err, common.BytesToStr(marshalJSON))
		return
	}

	obj = &dynamic{
		baseTransport: base.Transport{
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			DnsInfo: config.DnsInfo,
			BufferPool: sync.Pool{
				New: func() any {
					return make([]byte, 1024*32)
				},
			},
		},
	}

	if config.RandRule == "" {
		config.RandRule = rule.Timestamp
	}

	rand, err := rule.Factory(config.RandRule)
	if err != nil {
		return nil, err
	}

	obj.rule = rand

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
