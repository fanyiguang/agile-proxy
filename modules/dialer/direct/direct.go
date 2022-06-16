package direct

import (
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type direct struct {
	base.Dialer
}

func (d *direct) Dial(network string, host, port string) (conn net.Conn, err error) {
	return d.BaseDial(network, host, port)
}

func (d *direct) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	return d.BaseDialTimeout(network, host, port, timeout)
}

func New(jsonConfig json.RawMessage) (obj *direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "direct new")
		return
	}

	obj = &direct{
		Dialer: base.Dialer{
			Identity: plugin.Identity{
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
