package config_description

// https服务端配置
type httpsServer struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"` // 代理类型 https
	Name          string `json:"name"`
	TransportName string `json:"transport_name"`
	CrtPath       string `json:"crt_path"` // 证书路径（为空使用默认证书）
	KeyPath       string `json:"key_path"` // 密钥路径（为空使用默认密钥）
}

// https客户端配置
type httpsClient struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 https
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	CrtPath    string `json:"crt_path"`    // 证书路径（为空使用默认证书）
	KeyPath    string `json:"key_path"`    // 密钥路径（为空使用默认密钥）
	CaPath     string `json:"ca_path"`     // ca路径（为空使用系统ca）
	ServerName string `json:"server_name"` // 证书对应的server 使用默认证书的话server_name=localhost
	Mode       int    `json:"mode"`        // 转发模式 0-降级模式 1-严格模式
}

// https拨号器配置
type httpsDialer struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 https
	Name       string `json:"name"`
	CrtPath    string `json:"crt_path"`
	KeyPath    string `json:"key_path"`
	CaPath     string `json:"ca_path"`
	ServerName string `json:"server_name"`
}
