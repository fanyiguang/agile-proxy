package ha

import (
	"agile-proxy/model"
	"agile-proxy/modules/transport/base"
)

type Config struct {
	model.Base
	model.Identity
	model.PipelineInfos
	ClientNames string       `json:"client_names"` // 客户端名称已","隔开 例如：(ssh,ssl,socks5)
	DnsInfo     base.DnsInfo `json:"dns_info"`
}
