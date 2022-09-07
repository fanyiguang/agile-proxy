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
	model.Satellites
	RouteName string `json:"router_name"`
	KeyPath   string `json:"key_path"`
}
