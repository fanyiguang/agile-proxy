package socks5

import (
	"encoding/json"
	"github.com/pkg/errors"
	"nimble-proxy/modules/server"
)

type Socks5 struct {
	server.BaseServer
	Auth int
}

func (s Socks5) Start() (err error) {
	//TODO implement me
	panic("implement me")
}

func (s Socks5) Accept() {
	//TODO implement me
	panic("implement me")
}

func (s Socks5) Connect() {
	//TODO implement me
	panic("implement me")
}

func (s Socks5) Close() {
	//TODO implement me
	panic("implement me")
}

func New(strConfig string) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal([]byte(strConfig), &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Socks5{
		BaseServer: server.BaseServer{
			Ip:       config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
		Auth: config.Auth,
	}

	return
}
