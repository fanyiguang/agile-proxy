package direct

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

type Direct struct {
	base.Client
}

func (d *Direct) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	return d.Client.Dial(network, common.BytesToStr(host), common.BytesToStr(port))
}

func (d *Direct) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	return d.Client.DialTimeout(network, common.BytesToStr(host), common.BytesToStr(port), timeout)
}

func (d *Direct) createRoundTripper() (err error) {
	d.RoundTripper, err = d.CreateRoundTripper("", func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		deadline, ok := ctx.Deadline()
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return conn, err
		}

		if ok {
			now := time.Now()
			if deadline.After(now) {
				conn, err = d.Client.DialTimeout(network, host, port, deadline.Sub(now))
			} else {
				err = http.ErrHandlerTimeout
			}
		} else {
			conn, err = d.Client.Dial(network, host, port)
		}
		return
	})
	return
}

func (d *Direct) Run() (err error) {
	d.Client.Init()
	err = d.createRoundTripper()
	return
}

func (d *Direct) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *Direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Direct{
		Client: base.Client{
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			Satellites: config.Satellites,
			Mode:       config.Mode,
			DialerName: config.DialerName,
		},
	}

	return
}
