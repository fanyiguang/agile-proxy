package socks5

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	AuthMode  int    `json:"auth_mode"`
	Interface string `json:"interface"`
}
