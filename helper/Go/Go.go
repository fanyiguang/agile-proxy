package Go

import (
	"agile-proxy/helper/log"
	"fmt"
	"runtime/debug"
	"strings"
)

func Go(fun func()) {
	go func(f func()) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(strings.Repeat("+", 15))
				log.Error(fmt.Sprintf("%v %s", err, debug.Stack()))
			}
		}()

		f()
	}(fun)
}
