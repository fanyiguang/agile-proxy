package direct

import (
	"agile-proxy/model"
	"agile-proxy/modules/transport/base"
)

type Config struct {
	model.Base
	model.Identity
	model.PipelineInfos
	ClientName string       `json:"client_name"`
	DnsInfo    base.DnsInfo `json:"dns_info"`
}
