package https

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.Satellites
	DialerName string `json:"dialer_name"`
	CrtPath    string `json:"crt_path"`
	KeyPath    string `json:"key_path"`
	CaPath     string `json:"ca_path"`
	ServerName string `json:"server_name"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
