package direct

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/client"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport/base"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"sync"
)

type Direct struct {
	baseTransport base.Transport
	Client        client.Client // 传输器可以使用的客户端
}

func (d *Direct) Close() (err error) {
	return nil
}

func (d *Direct) Transport(cConn net.Conn, host, port []byte) (err error) {
	if d.Client != nil {
		host, err = d.baseTransport.GetHost(host)
		if err != nil {
			return
		}

		var sConn net.Conn
		sConn, err = d.Client.Dial("tcp", host, port)
		if err != nil {
			return
		}

		defer sConn.Close()
		d.baseTransport.AsyncSendMsg(fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		d.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("Client is nil")
	}
	return
}

func New(jsonConfig json.RawMessage) (obj *Direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "direct new")
		return
	}

	obj = &Direct{
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

	if config.ClientName != "" {
		obj.Client = client.GetClient(config.ClientName)
	}
	return
}
