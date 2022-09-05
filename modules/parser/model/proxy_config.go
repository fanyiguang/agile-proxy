package model

import "encoding/json"

type ProxyConfig struct {
	ClientConfig []json.RawMessage `json:"client"`
	ServerConfig []json.RawMessage `json:"server"`
	DialerConfig []json.RawMessage `json:"dialer"`
	RouteConfig  []json.RawMessage `json:"router"`
	MsgConfig    []json.RawMessage `json:"msg"`
	LogPath      string            `json:"log_path"`
	LogLevel     string            `json:"log_level"`
}
