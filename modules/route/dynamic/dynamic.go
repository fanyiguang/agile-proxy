package dynamic

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client"
	"agile-proxy/modules/route/base"
	"agile-proxy/modules/route/dynamic/rule"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
)

type dynamic struct {
	baseTransport base.Transport
	clients       []client.Client // 动态类型的传输器客户端可以为多个
	rule          rule.Rule
	randRule      string
	clientNames   string
	clientsLen    int
}

func (d *dynamic) Transport(cConn net.Conn, host, port []byte) (err error) {
	if d.clients != nil {
		host, err = d.baseTransport.Dns.GetHost(host)
		if err != nil {
			return
		}

		var sConn net.Conn
		sConn, err = d.clients[d.getClientIndex()].Dial("tcp", host, port)
		if err != nil {
			return
		}

		defer sConn.Close()
		d.baseTransport.AsyncSendMsg(d.baseTransport.Name(), -1, fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		d.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("Client is nil")
	}
	return
}

func (d *dynamic) Run() (err error) {
	err = d.init()
	return
}

func (d *dynamic) Close() (err error) {
	return
}

func (d *dynamic) getClientIndex() (idx int) {
	return d.rule.Intn(d.clientsLen)
}

func (d *dynamic) init() (err error) {
	d.baseTransport.Init()
	if d.randRule == "" {
		d.randRule = rule.Timestamp
	}

	d.rule, err = rule.Factory(d.randRule)
	if err != nil {
		return
	}

	if d.clientNames != "" {
		clientNames := strings.Split(d.clientNames, ",")
		for _, clientName := range clientNames {
			_client := client.GetClient(clientName)
			if _client != nil {
				d.clients = append(d.clients, _client)
			}
		}
	}

	d.clientsLen = len(d.clients)
	return
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
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			Dns:           assembly.CreateDns(config.DnsInfo.Server, config.DnsInfo.LocalDns),
			BufferPool:    common.CreateByteBufferSyncPool(1024 * 32),
			PipelineInfos: config.PipelineInfos,
		},
		randRule:    config.RandRule,
		clientNames: config.ClientNames,
	}

	return
}
