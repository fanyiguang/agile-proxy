package direct

import (
	"agile-proxy/modules/transport/base"
)

type Config struct {
	Type       string       `json:"type"`
	Name       string       `json:"name"`
	ClientName string       `json:"client_name"`
	DnsInfo    base.DnsInfo `json:"dns_info"`
}
