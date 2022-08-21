package direct

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client"
	"agile-proxy/modules/transport/base"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
)

type direct struct {
	baseTransport base.Transport
	Client        client.Client // 传输器可以使用的客户端
	clientName    string
}

func (d *direct) Run() (err error) {
	d.init()
	return
}

func (d *direct) Close() (err error) {
	return
}

func (d *direct) Transport(cConn net.Conn, host, port []byte) (err error) {
	if d.Client != nil {
		host, err = d.baseTransport.Dns.GetHost(host)
		if err != nil {
			return
		}

		var sConn net.Conn
		sConn, err = d.Client.Dial("tcp", host, port)
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

func (d *direct) init() {
	d.baseTransport.Init()
	if d.clientName != "" {
		d.Client = client.GetClient(d.clientName)
	}
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
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			Dns:           assembly.CreateDns(config.DnsInfo.Server, config.DnsInfo.LocalDns),
			BufferPool:    common.CreateByteBufferSyncPool(1024 * 32),
			PipelineInfos: config.PipelineInfos,
		},
		clientName: config.ClientName,
	}

	return
}
