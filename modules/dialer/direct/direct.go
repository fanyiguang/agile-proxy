package direct

import (
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Direct struct {
	base.Dialer
}

func (d *Direct) Dial(network string, host, port string) (conn net.Conn, err error) {
	return d.BaseDial(network, host, port)
}

func (d *Direct) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	return d.BaseDialTimeout(network, host, port, timeout)
}

func New(jsonConfig json.RawMessage) (obj *Direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "direct new")
		return
	}

	obj = &Direct{
		Dialer: base.Dialer{
			IdentInfo: plugin.IdentInfo{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			IFace: config.Interface,
		},
	}

	return
}
