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
	"sync"
	"time"
)

type Transport struct {
	plugin.Identity
	OutMsg     plugin.PipelineOutput
	DnsInfo    model.DnsInfo
	BufferPool sync.Pool
}

func (t *Transport) AsyncSendMsgToIpc(msg string) {
	// 异步对外发送消息，减少对主流程的影响
	// 对外保持0信任原则，设置超时时间如果
	// 外部阻塞也不会导致协程泄漏。
	Go.Go(func() {
		select {
		case t.OutMsg.Ch <- plugin.OutputMsg{
			Content:    msg,
			ModuleName: t.Name(),
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

func (t *Transport) Copy(sConn net.Conn, cConn net.Conn) {
	errCh := make(chan error, 1)
	go t.copyBuffer(sConn, cConn, errCh)
	go t.copyBuffer(cConn, sConn, errCh)
	t.wait(errCh)
}

func (t *Transport) copyBuffer(det net.Conn, src net.Conn, errCh chan error) {
	buffer := t.BufferPool.Get().([]byte)
	defer t.BufferPool.Put(buffer)
	_, err := io.CopyBuffer(det, src, buffer)
	errCh <- err
}

func (t *Transport) wait(errCh chan error) {
	err := <-errCh
	if err != nil && err != io.EOF {
		log.WarnF("wait failed：%v", err)
	}
}
