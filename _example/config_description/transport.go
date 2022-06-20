package config_description

// 直连传输器配置
type directTransport struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	ClientName string `json:"client_name"`
	DnsInfo    struct {
		Server   string `json:"server"`    // dns server 地址 114.114.114.114 or 8.8.8.8
		LocalDns bool   `json:"local_dns"` //是否开启本地dns true-开启 false-关闭
	} `json:"dns_info"`
}

// 动态传输器配置
type dynamicTransport struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	ClientNames string `json:"client_names"` // 客户端名称已","隔开 例如：(ssh,ssl,socks5)
	RandRule    string `json:"rand_rule"`    // default timestamp
	DnsInfo     struct {
		Server   string `json:"server"`    // dns server 地址 114.114.114.114 or 8.8.8.8
		LocalDns bool   `json:"local_dns"` //是否开启本地dns true-开启 false-关闭
	} `json:"dns_info"`
}

// 高可用传输器配置
type haTransport struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	ClientNames string `json:"client_names"` // 客户端名称已","隔开 例如：(ssh,ssl,socks5)
	DnsInfo     struct {
		Server   string `json:"server"`    // dns server 地址 114.114.114.114 or 8.8.8.8
		LocalDns bool   `json:"local_dns"` //是否开启本地dns true-开启 false-关闭
	} `json:"dns_info"`
}
