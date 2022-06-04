package socks5

type Config struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	AuthMode int    `json:"auth_mode"`
}
