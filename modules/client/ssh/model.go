package ssh

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	DialerName string `json:"dialer_name"`
	KeyPath    string `json:"key_path"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
