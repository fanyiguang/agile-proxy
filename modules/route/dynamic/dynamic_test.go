package dynamic

import (
	"agile-proxy/modules/route/base"
	"agile-proxy/modules/route/dynamic/rule"
	"fmt"
	"testing"
)

func TestGetClientIndex(t *testing.T) {
	metas := []struct {
		clientLen int
		rangeLoop int
	}{
		{10, 300},
		{20, 100},
		{10, 600},
		{30, 100},
		{60, 1000},
	}

	for _, meta := range metas {
		d := dynamic{
			baseTransport: base.Transport{},
			clientsLen:    meta.clientLen,
		}

		d.rule, _ = rule.Factory(rule.Timestamp)
		for i := 0; i < meta.rangeLoop; i++ {
			index := d.getClientIndex()
			if index >= d.clientsLen {
				t.Error("out of range")
				return
			}
		}
	}
	fmt.Println("success")
}
