package assembly

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/dns"
	"agile-proxy/helper/log"
	helperNet "agile-proxy/helper/net"
	"github.com/pkg/errors"
	"net"
)

type Dns struct {
	Server   string
	LocalDns bool
}

func (d *Dns) GetHost(host []byte) (newHost []byte, err error) {
	// 增加net.ParseIP校验的目的是防止客户端socks5协议
	// 告诉我们是域名类型的host但host却是字符串的ip时去
	// 走dns的情况。
	if !d.LocalDns || net.ParseIP(common.BytesToStr(host)) != nil { // 没有设置本地dns
		newHost = host
		return
	}

	var strHost, port string
	strHost, port, err = helperNet.SplitHostAndPort(common.BytesToStr(host))
	if err != nil {
		return
	}

	if d.Server != "" { // 配置文件的dns服务器不为空
		lookupHost, err := dns.LookupHost(strHost, d.Server)
		if err == nil {
			return common.StrToBytes(helperNet.JoinHostAndPort(lookupHost[0], port)), nil
		}

		log.WarnF("dns.LookupHost failed: %v, host: %v, dns_server: %v", err, strHost, d.Server)
	}

	// 走系统配置的dns
	// TODO ips为空时 _err等于nil吗？
	ips, _err := net.LookupIP(strHost)
	if _err != nil {
		err = errors.Wrap(_err, "net.LookupIP")
		return
	}
	for _, ip := range ips {
		return common.StrToBytes(helperNet.JoinHostAndPort(ip.String(), port)), nil
	}

	err = errors.Wrap(errors.New("ips len is 0"), "")
	return
}

func CreateDns(server string, localDns bool) Dns {
	return Dns{
		Server:   server,
		LocalDns: localDns,
	}
}
