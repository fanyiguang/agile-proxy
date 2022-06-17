package base

type DnsInfo struct {
	Server   string `json:"server"`
	LocalDns bool   `json:"local_dns"`
}
