package base

import (
	"io"
	"net"
	"time"
)

type Transport struct {
	TransportType string
	TransportName string
	ClientName    string
}

func (t *Transport) MyCopy(sConn net.Conn, cConn net.Conn) {
	cDoneCh, sDoneCh := make(chan struct{}), make(chan struct{})
	go t.myCopyN(sConn, cConn, cDoneCh)
	go t.myCopyN(cConn, sConn, sDoneCh)
	t.wait(cDoneCh, sDoneCh)
}

func (t *Transport) myCopyN(det net.Conn, src net.Conn, done chan struct{}) {
	for {
		_ = src.SetReadDeadline(time.Now().Add(2 * time.Minute))
		_ = det.SetWriteDeadline(time.Now().Add(2 * time.Minute))
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
