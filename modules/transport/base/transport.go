package base

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/log"
	"agile-proxy/modules/plugin"
	"io"
	"net"
	"sync"
	"time"
)

type Transport struct {
	plugin.Identity
	OutMsg     plugin.PipelineOutput
	Dns        plugin.Dns
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
