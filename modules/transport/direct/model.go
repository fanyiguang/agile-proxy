package direct

import "nimble-proxy/modules/transport"

type Config struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	ClientName string            `json:"client_name"`
	DnsInfo    transport.DnsInfo `json:"dns_info"`
}
