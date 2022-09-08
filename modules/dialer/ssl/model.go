package ssl

import "agile-proxy/model"

type Config struct {
	model.Base
	model.Net
	model.Identity
	model.Satellites
	CrtPath    string `json:"crt_path"`
	KeyPath    string `json:"key_path"`
	CaPath     string `json:"ca_path"`
	ServerName string `json:"server_name"`
	Interface  string `json:"interface"`
	AuthMode   int    `json:"auth_mode"`
}
