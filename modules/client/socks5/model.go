package socks5

type Config struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"`
	DialerName string `json:"dialer_name"`
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"`
}
