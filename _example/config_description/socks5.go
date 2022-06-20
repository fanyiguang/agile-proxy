package config_description

// socks5服务端配置
type socks5Server struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"` // 代理类型 socks5
	Name          string `json:"name"`
	TransportName string `json:"transport_name"` // 传输器名称
	AuthMode      int    `json:"auth_mode"`      // 认证模式 0-不允许匿名模式 1-允许匿名模式
}

// socks5客户端配置
type socks5Client struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 socks5
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"` // 拨号器名称
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}

// socks5拨号器配置
type socks5Dialer struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"` // 代理类型 socks5
	Name     string `json:"name"`
	AuthMode int    `json:"auth_mode"`
}
