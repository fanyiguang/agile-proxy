package direct

import (
	"agile-proxy/modules/transport/base"
	"agile-proxy/modules/transport/model"
	"fmt"
	"testing"
)

func TestGetHost(t *testing.T) {
	metas := []struct {
		server   string
		localDns bool
		host     string
		res      string
	}{
		{"", false, "www.baidu.com", "www.baidu.com"},
		{"114.114.114.114", false, "www.baidu.com", "www.baidu.com"},
		//{"114.114.114.114:53", true, "www.amazon.com", "99.84.235.200"},
		//{"", true, "www.amazon.com", "99.84.235.200"},
		{"114.114.114.114", true, "99.84.235.200", "99.84.235.200"},
		{"", true, "99.84.235.200", "99.84.235.200"},
		{"", false, "99.84.235.200", "99.84.235.200"},
	}

	for _, meta := range metas {
		direct := direct{
			baseTransport: base.Transport{
				DnsInfo: model.DnsInfo{
					Server:   meta.server,
					LocalDns: meta.localDns,
				},
			},
		}

		host, err := direct.baseTransport.GetHost([]byte(meta.host))
		if err != nil {
			t.Error(err)
			return
		}

		if string(host) != meta.res {
			t.Error(err)
			return
		}
	}

	fmt.Println("success")

}
