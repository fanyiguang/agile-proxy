package base

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/dns"
	"agile-proxy/helper/log"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport/model"
	"github.com/pkg/errors"
	"io"
	"net"
	"time"
)

type Transport struct {
	plugin.IdentInfo
	OutMsg  plugin.PipelineOutput
	DnsInfo model.DnsInfo
}

func (t *Transport) AsyncSendMsg(msg string) {
	// 异步对外发送消息，减少对主流程的影响
	// 对外保持0信任原则，设置超时时间如果
	// 外部阻塞就不会协程泄漏。
	Go.Go(func() {
		select {
		case t.OutMsg.Ch <- plugin.OutputMsg{
			Content: msg,
			Module:  t.Type(),
		}:
		case <-time.After(time.Second):
			log.InfoF("pipeline message lock: %v %v", msg, t.Name())
		}
	})
}

func (t *Transport) GetHost(host []byte) (newHost []byte, err error) {
	// 增加net.ParseIP校验的目的是防止客户端socks5协议
	// 告诉我们是域名类型的host但host却是字符串的ip时去
	// 走dns的情况。
	if !t.DnsInfo.LocalDns || net.ParseIP(common.BytesToStr(host)) != nil { // 没有设置本地dns
		newHost = host
		return
	}

	strHost := common.BytesToStr(host)
	if t.DnsInfo.Server != "" { // 配置文件的dns服务器不为空
		lookupHost, err := dns.LookupHost(strHost, t.DnsInfo.Server)
		if err == nil {
			return common.StrToBytes(lookupHost[0]), nil
		}

		log.WarnF("dns.LookupHost failed: %v, host: %v, dns_server: %v", err, strHost, t.DnsInfo.Server)
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

func (t *Transport) MyCopy(sConn net.Conn, cConn net.Conn) {
	cDoneCh, sDoneCh := make(chan struct{}), make(chan struct{})
	go t.myCopyN(sConn, cConn, cDoneCh)
	go t.myCopyN(cConn, sConn, sDoneCh)
	t.wait(cDoneCh, sDoneCh)
}

func (t *Transport) myCopyN(det net.Conn, src net.Conn, done chan struct{}) {
	for {
		_ = src.SetReadDeadline(time.Now().Add(time.Minute))
		_ = det.SetWriteDeadline(time.Now().Add(time.Minute))
		_, err := io.CopyN(det, src, 1024*32)
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
	}
	close(done)
}

func (t *Transport) wait(cDoneCh, sDoneCh chan struct{}) {
	<-cDoneCh
	<-sDoneCh
}
