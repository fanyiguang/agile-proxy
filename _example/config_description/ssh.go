package config_description

// ssh服务端配置
type sshServer struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"` // 代理类型 ssh
	Name          string `json:"name"`
	TransportName string `json:"transport_name"`
	KeyPath       string `json:"key_path"` // 公钥路径
}

// ssh客户端配置
type sshClient struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"` // 代理类型 ssh
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	KeyPath    string `json:"key_path"` // 私钥路径
	Mode       int    `json:"mode"`     // 转发模式 0-降级模式 1-严格模式
}

// ssh拨号器配置
type sshDialer struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"` // 代理类型 ssh
	Name     string `json:"name"`
	KeyPath  string `json:"key_path"` // 私钥路径
}
