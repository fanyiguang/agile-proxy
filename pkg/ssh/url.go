package ssh

func HttpDefaultPort(scheme string) (port string) {
	switch scheme {
	case "http":
		port = "80"
	case "https":
		port = "443"
	}
	return
}
