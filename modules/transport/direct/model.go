package direct

import (
	"agile-proxy/modules/transport/model"
)

type Config struct {
	Type       string        `json:"type"`
	Name       string        `json:"name"`
	ClientName string        `json:"client_name"`
	DnsInfo    model.DnsInfo `json:"dns_info"`
}
