package ssl

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	DialerName string `json:"dialer_name"`
	CrtPath    string `json:"crt_path"`
	KeyPath    string `json:"key_path"`
	CaPath     string `json:"ca_path"`
	ServerName string `json:"server_name"`
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
