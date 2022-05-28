package common

func CloseChan[a any](ch chan a) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
