package socks5

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	DialerName string `json:"dialer_name"`
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
