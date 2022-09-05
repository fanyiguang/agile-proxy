package https

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	RouteName string `json:"route_name"`
	CrtPath   string `json:"crt_path"`
	KeyPath   string `json:"key_path"`
}
