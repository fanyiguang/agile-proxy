package ssh

import "agile-proxy/model"

type DirectForward struct {
	DesAddr string
	DesPort uint32

	OriginAddr string
	OriginPort uint32
}

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.PipelineInfos
	TransportName string `json:"transport_name"`
	KeyPath       string `json:"key_path"`
}
