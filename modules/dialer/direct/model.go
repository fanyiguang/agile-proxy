package direct

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	Interface string `json:"interface"`
}
