package config_description

// http服务端配置
type httpServer struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"` // 代理类型 http
	Name          string `json:"name"`
	TransportName string `json:"transport_name"`
}

// http客户端配置
type httpClient struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 http
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}

// http拨号器配置
type httpDialer struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"` // 代理类型 http
	Name     string `json:"name"`
}
