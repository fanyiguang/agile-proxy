package server

var servers []Server

func GetAllServer() []Server {
	return servers
}

func CloseAllServers() {
	for _, server := range servers {
		if server != nil {
			_ = server.Close()
		}
	}
}
