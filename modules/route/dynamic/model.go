package dynamic

import (
	"agile-proxy/model"
	"agile-proxy/modules/route/base"
)

type Config struct {
	model.Base
	model.Identity
	model.Satellites
	ClientNames string       `json:"client_names"` // 客户端名称已","隔开 例如：(ssh,ssl,socks5)
	RandRule    string       `json:"rand_rule"`    // default timestamp
	DnsInfo     base.DnsInfo `json:"dns_info"`
}
