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

type direct struct {
	baseTransport base.Transport
	Client        client.Client // 传输器可以使用的客户端
}

func (d *direct) Close() (err error) {
	if d.Client != nil {
		err = d.Client.Close()
	}
	return
}

func (d *direct) Transport(cConn net.Conn, host, port []byte) (err error) {
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
		d.baseTransport.AsyncSendMsgToIpc(fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		d.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("Client is nil")
	}
	return
}

func New(jsonConfig json.RawMessage) (obj *direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		marshalJSON, _ := jsonConfig.MarshalJSON()
		err = errors.Wrap(err, common.BytesToStr(marshalJSON))
		return
	}

	obj = &direct{
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
