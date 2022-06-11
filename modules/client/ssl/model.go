package ssl

type Config struct {
	Ip         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	CrtPath    string `json:"crt_path"`
	KeyPath    string `json:"key_path"`
	CaPath     string `json:"ca_path"`
	AuthMode   int    `json:"auth_mode"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
