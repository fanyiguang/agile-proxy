package parser

type ProxyConfigInfo struct {
	ClientConfig    []string `json:"client_config"`
	ServerConfig    []string `json:"server_config"`
	DialerConfig    []string `json:"dialer_config"`
	TransportConfig []string `json:"transport_config"`
	LogPath         string   `json:"log_path"`
	LogLevel        string   `json:"log_level"`
}
