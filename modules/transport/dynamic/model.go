package dynamic

import (
	"agile-proxy/modules/transport/model"
)

type Config struct {
	Type        string        `json:"type"`
	Name        string        `json:"name"`
	ClientNames string        `json:"client_names"` // 客户端名称已","隔开 例如：(ssh,ssl,socks5)
	RandRule    string        `json:"rand_rule"`    // default timestamp
	DnsInfo     model.DnsInfo `json:"dns_info"`
}
