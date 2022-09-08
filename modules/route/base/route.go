package base

import (
	"agile-proxy/helper/log"
	helperNet "agile-proxy/helper/net"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/satellite"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Transport struct {
	assembly.Dns
	assembly.Identity
	assembly.Pipeline
	model.Satellites
	BufferPool sync.Pool
}

func (t *Transport) Copy(sConn net.Conn, cConn net.Conn) {
	errCh := make(chan error, 1)
	go t.copyBuffer(sConn, cConn, errCh)
	go t.copyBuffer(cConn, sConn, errCh)
	t.wait(errCh)
}

func (t *Transport) HttpCopy(w http.ResponseWriter, resp *http.Response) {
	helperNet.CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	errCh := make(chan error, 1)
	t.copyBuffer(w, resp.Body, errCh)
}

func (t *Transport) copyBuffer(dst io.Writer, src io.Reader, errCh chan error) {
	buffer := t.BufferPool.Get().([]byte)
	defer t.BufferPool.Put(buffer)
	_, err := io.CopyBuffer(dst, src, buffer)
	errCh <- err
}

func (t *Transport) Init() {
	if !strings.Contains(t.Dns.Server, ":") {
		t.Dns.Server = net.JoinHostPort(t.Dns.Server, "53")
	}

	for _, _satellite := range t.Satellites.Satellites {
		_msg := satellite.GetSatellite(_satellite.Name)
		if _msg != nil {
			msgPipeline, level := _msg.Subscribe(t.Name(), t.Pipeline.PipeCh, _satellite.Level)
			t.Subscribe(_satellite.Name, msgPipeline, level)
		} else {
			log.WarnF("%v transport get msg failed pipeline name: %v", t.Name(), _satellite.Name)
		}
	}
}

func (t *Transport) wait(errCh chan error) {
	err := <-errCh
	if err != nil && err != io.EOF {
		log.WarnF("wait failed：%v", err)
	}
}
