package direct

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/transport/base"
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
		{"114.114.114.114:53", true, "www.classical.com", "34.116.240.52"},
		{"", true, "www.classical.com", "34.116.240.52"},
		{"114.114.114.114", true, "99.84.235.200", "99.84.235.200"},
		{"", true, "99.84.235.200", "99.84.235.200"},
		{"", false, "99.84.235.200", "99.84.235.200"},
	}

	for _, meta := range metas {
		direct := direct{
			baseTransport: base.Transport{
				Dns: assembly.Dns{
					Server:   meta.server,
					LocalDns: meta.localDns,
				},
			},
		}

		host, err := direct.baseTransport.Dns.GetHost([]byte(meta.host))
		if err != nil {
			t.Error(err)
			return
		}

		if string(host) != meta.res {
			t.Error(fmt.Sprintf("host not eq meta res %v %v", string(host), meta.res))
			return
		}
	}

	fmt.Println("success")

}
