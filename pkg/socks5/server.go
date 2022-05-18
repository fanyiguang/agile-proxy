package socks5

func IsSocks5(t byte) bool {
	if t == tag {
		return true
	} else {
		return false
	}
}
