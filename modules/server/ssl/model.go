package ssl

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	RouteName string `json:"router_name"`
	CrtPath   string `json:"crt_path"`
	KeyPath   string `json:"key_path"`
	AuthMode  int    `json:"auth_mode"` // 认证模式 0-允许匿名模式 1-不允许匿名模式
}
