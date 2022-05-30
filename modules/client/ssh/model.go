package ssh

type Config struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	RsaPath    string `json:"rsa_path"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
