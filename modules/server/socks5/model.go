package socks5

type Config struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"`
	TransportName string `json:"transport_name"`
	AuthMode      int    `json:"auth_mode"`      // 认证模式 0-允许匿名模式 1-不允许匿名模式
	TransportMode int    `json:"transport_mode"` // 转发模式 0-降级模式 1-严格模式
}
