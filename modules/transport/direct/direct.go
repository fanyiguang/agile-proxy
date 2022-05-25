package direct

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/common"
	"nimble-proxy/helper/dns"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client"
	"nimble-proxy/modules/transport/base"
	"nimble-proxy/modules/transport/model"
)

type Direct struct {
	baseTransport base.Transport
	Client        client.Client // 传输器可以使用的客户端
	DnsInfo       model.DnsInfo
}

func (d *Direct) Close() (err error) {
	return nil
}

func (d *Direct) Transport(cConn net.Conn, host, port []byte) (err error) {
	if d.Client != nil {
		host, err = d.getHost(host)
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

func (d *Direct) getHost(host []byte) (newHost []byte, err error) {
	// 增加net.ParseIP校验的目的是防止客户端socks5协议
	// 告诉我们是域名类型的host但host却是字符串的ip时去
	// 走dns的情况。
	if !d.DnsInfo.LocalDns || net.ParseIP(common.BytesToStr(host)) != nil { // 没有设置本地dns
		newHost = host
		return
	}

	strHost := common.BytesToStr(host)
	if d.DnsInfo.Server != "" { // 配置文件的dns服务器不为空
		lookupHost, err := dns.LookupHost(strHost, d.DnsInfo.Server)
		if err == nil {
			return common.StrToBytes(lookupHost[0]), nil
		}

		log.WarnF("dns.LookupHost failed: %v, host: %v, dns_server: %v", err, strHost, d.DnsInfo.Server)
	}

	// 走系统配置的dns
	// TODO ips为空时 _err等于nil吗？
	ips, _err := net.LookupIP(strHost)
	if _err != nil {
		err = errors.Wrap(_err, "net.LookupIP")
		return
	}
	for _, ip := range ips {
		return common.StrToBytes(ip.String()), nil
	}

	err = errors.Wrap(errors.New("ips len is 0"), "")
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
			ClientName:    config.ClientName,
		},
		DnsInfo: config.DnsInfo,
	}

	if config.ClientName != "" {
		obj.Client = client.GetClient(config.ClientName)
	}
	return
}
