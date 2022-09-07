package direct

import (
	"agile-proxy/model"
	"agile-proxy/modules/route/base"
)

type Config struct {
	model.Base
	model.Identity
	model.Satellites
	ClientName string       `json:"client_name"`
	DnsInfo    base.DnsInfo `json:"dns_info"`
}
