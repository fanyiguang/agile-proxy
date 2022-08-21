package base

import (
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/msg"
	"io"
	"net"
	"strings"
	"sync"
)

type Transport struct {
	assembly.Dns
	assembly.Identity
	assembly.Pipeline
	model.PipelineInfos
	BufferPool sync.Pool
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

func (t *Transport) Init() {
	if !strings.Contains(t.Dns.Server, ":") {
		t.Dns.Server = net.JoinHostPort(t.Dns.Server, "53")
	}

	for _, pipelineInfo := range t.PipelineInfo {
		_msg := msg.GetMsg(pipelineInfo.Name)
		if _msg != nil {
			msgPipeline, level := _msg.Subscribe(t.Name(), t.Pipeline.PipeCh, pipelineInfo.Level)
			t.Subscribe(pipelineInfo.Name, msgPipeline, level)
		} else {
			log.WarnF("%v transport get msg failed pipeline name: %v", t.Name(), pipelineInfo.Name)
		}
	}
}

func (t *Transport) wait(errCh chan error) {
	err := <-errCh
	if err != nil && err != io.EOF {
		log.WarnF("wait failed：%v", err)
	}
}
