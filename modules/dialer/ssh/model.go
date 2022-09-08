package ssh

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.Satellites
	KeyPath   string `json:"key_path"`
	Interface string `json:"interface"`
}
