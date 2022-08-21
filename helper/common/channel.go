package common

import "time"

func CloseChan[a any](ch chan a) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}

func ReliableChanSend[a any](ch chan a, val a) {
	select {
	case ch <- val:
	default:
	}
}

func SendChanTimeout[a any](ch chan a, val a, timeout int) {
	select {
	case ch <- val:
	case <-time.After(time.Duration(timeout) * time.Second):
	}
}

func AcceptChanTimeout[a any](ch chan a, timeout int) (val a, ok bool) {
	select {
	case val = <-ch:
		ok = true
	case <-time.After(time.Duration(timeout) * time.Second):
		ok = false
	}
	return
}
