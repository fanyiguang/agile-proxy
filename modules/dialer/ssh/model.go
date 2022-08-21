package ssh

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	KeyPath   string `json:"key_path"`
	Interface string `json:"interface"`
}
