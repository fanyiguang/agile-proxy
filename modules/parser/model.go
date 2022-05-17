package parser

type ProxyConfigInfo struct {
	ClientConfig    interface{} `json:"client_config"`
	ServerConfig    interface{} `json:"server_config"`
	DialerConfig    interface{} `json:"dialer_config"`
	TransportConfig interface{} `json:"transport_config"`
	LogPath         string      `json:"log_path"`
	LogLevel        string      `json:"log_level"`
}
