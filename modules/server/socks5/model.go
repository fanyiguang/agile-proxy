package socks5

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.Satellites
	RouteName string `json:"router_name"`
	AuthMode  int    `json:"auth_mode"` // 认证模式 0-不允许匿名模式 1-允许匿名模式
}
