package direct

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Identity
	model.Satellites
	DialerName string `json:"dialer_name"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
