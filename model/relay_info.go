package model

import "encoding/json"

type RelayInfo struct {
	Ip             string      `json:"ip"`
	Region         string      `json:"region"`
	WhiteList      []string    `json:"white_list"`
	BlackList      []string    `json:"black_list"`
	CloudRegionKey json.Number `json:"cloud_region_key"`
	SecretKey      string      `json:"secret_key"`
	Country        string      `json:"country"`
	UserID         json.Number `json:"user_id"`
	MachineString  string      `json:"machine_string"`
}

type ProxyIpLocation struct {
	Country string `json:"country"`
}
