package direct

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type direct struct {
	base.Dialer
}

func (d *direct) Run() (err error) {
	return
}

func (d *direct) Close() (err error) {
	return
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
			Net:        assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			Satellites: config.Satellites,
			IFace:      config.Interface,
		},
	}

	return
}
