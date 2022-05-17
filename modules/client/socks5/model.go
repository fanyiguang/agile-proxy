package socks5

type Config struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"`
	TransportName string `json:"transport_name"`
	Auth          int    `json:"auth"`
}
