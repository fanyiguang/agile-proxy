package Go

import (
	"fmt"
	"nimble-proxy/helper/log"
	"runtime/debug"
	"strings"
)

func Go(fun func()) {
	go func(f func()) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(strings.Repeat("+", 15))
				log.Error(fmt.Sprintf("%s", debug.Stack()))
			}
		}()

		f()
	}(fun)
}
