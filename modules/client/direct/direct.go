package direct

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Direct struct {
	base.Client
}

func (d *Direct) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	if d.Dialer != nil {
		conn, err = d.Dialer.Dial(network, common.BytesToStr(host), d.GetStrPort(port))
		if err == nil || d.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("d.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(common.BytesToStr(host), d.GetStrPort(port)))
	if err != nil {
		err = errors.Wrap(err, "direct Dial")
	}
	return
}

func (d *Direct) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	if d.Dialer != nil {
		conn, err = d.Dialer.DialTimeout(network, common.BytesToStr(host), d.GetStrPort(port), timeout)
		if err == nil || d.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("d.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(common.BytesToStr(host), d.GetStrPort(port)))
	if err != nil {
		err = errors.Wrap(err, "direct Dial")
	}
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
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Mode: config.Mode,
		},
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	return
}
