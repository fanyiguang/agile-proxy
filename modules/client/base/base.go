package base

type Client struct {
	Host     string
	Port     string
	Username string
	Password string
	Mode     int // 0-降级模式（如果有配置连接器且连接器无法使用会走默认网络，默认为降级模式） 1-严格模式（如果有配置连接器且连接器无法使用则直接返回失败）
}
