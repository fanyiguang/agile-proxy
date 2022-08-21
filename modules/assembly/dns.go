package assembly

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/dns"
	"agile-proxy/helper/log"
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

	strHost := common.BytesToStr(host)
	if d.Server != "" { // 配置文件的dns服务器不为空
		lookupHost, err := dns.LookupHost(strHost, d.Server)
		if err == nil {
			return common.StrToBytes(lookupHost[0]), nil
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
		return common.StrToBytes(ip.String()), nil
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
