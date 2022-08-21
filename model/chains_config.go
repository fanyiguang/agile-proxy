package model

type ChainsConfig struct {
	Type   int    `json:"type"`
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	Scheme string `json:"scheme"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Params string `json:"params"`
}
