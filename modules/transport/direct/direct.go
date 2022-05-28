package direct

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/modules/client"
	"nimble-proxy/modules/transport/base"
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

		d.baseTransport.MyCopy(sConn, cConn)
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
			TransportType: config.Type,
			TransportName: config.Name,
			DnsInfo:       config.DnsInfo,
		},
	}

	if config.ClientName != "" {
		obj.Client = client.GetClient(config.ClientName)
	}
	return
}
