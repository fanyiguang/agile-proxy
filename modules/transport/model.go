package transport

type BaseTransport struct {
	Type       string
	Name       string
	ClientName string
}

type DnsInfo struct {
	Server   string `json:"server"`
	LocalDns bool   `json:"local_dns"`
}
