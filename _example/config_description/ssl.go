package config_description

// ssl服务端配置
type sslServer struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"` // 代理类型 ssl
	Name          string `json:"name"`
	TransportName string `json:"transport_name"`
	CrtPath       string `json:"crt_path"`  // 证书路径（为空使用默认证书）
	KeyPath       string `json:"key_path"`  // 密钥路径（为空使用默认密钥）
	AuthMode      int    `json:"auth_mode"` // 认证模式 0-允许匿名模式 1-不允许匿名模式
}

// ssl客户端配置
type sslClient struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 ssl
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	CrtPath    string `json:"crt_path"`    // 证书路径（为空使用默认证书）
	KeyPath    string `json:"key_path"`    // 密钥路径（为空使用默认密钥）
	CaPath     string `json:"ca_path"`     // ca路径（为空使用系统ca）
	ServerName string `json:"server_name"` // 证书对应的server 使用默认证书的话server_name=localhost
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}

// ssl拨号器配置
type sslDialer struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 ssl
	Name       string `json:"name"`
	CrtPath    string `json:"crt_path"`    // 证书路径（为空使用默认证书）
	KeyPath    string `json:"key_path"`    // 密钥路径（为空使用默认密钥）
	CaPath     string `json:"ca_path"`     // ca路径（为空使用系统ca）
	ServerName string `json:"server_name"` // 证书对应的server 使用默认证书的话server_name=localhost
	AuthMode   int    `json:"auth_mode"`
}
