package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// LookupHost
// ips为空时err!=nil
func LookupHost(domain, server string) (ips []string, err error) {
	var msg *dns.Msg
	client := dns.Client{}
	dnsMsg := dns.Msg{}
	dnsMsg.SetQuestion(fmt.Sprintf("%s.", domain), dns.TypeA)
	msg, _, err = client.Exchange(&dnsMsg, server)
	if err != nil {
		err = errors.Wrap(err, "client.Exchange")
		return
	}

	for _, d := range msg.Answer {
		switch d.Header().Rrtype {
		case dns.TypeA:
			ips = append(ips, d.(*dns.A).A.String())
		case dns.TypeAAAA:
			ips = append(ips, d.(*dns.AAAA).AAAA.String())
		}
	}

	if len(ips) < 1 {
		err = errors.New("ips is nil")
	}
	return
}
