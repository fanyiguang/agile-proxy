package rule

import (
	"agile-proxy/modules/route/dynamic/rule/timestamp"
	"fmt"
	"github.com/pkg/errors"
)

type Rule interface {
	Int() int
	Intn(n int) int
}

func Factory(t string) (rand Rule, err error) {
	switch t {
	case Timestamp:
		rand, err = timestamp.New()
	default:
		err = errors.New(fmt.Sprintf("random type invalid %v", t))
	}
	return
}
