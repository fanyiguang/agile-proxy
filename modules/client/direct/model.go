package direct

type Config struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	DialerName string `json:"dialer_name"`
	Mode       int    `json:"mode"` // 转发模式 0-降级模式 1-严格模式
}
