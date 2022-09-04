package server

var servers = make(map[string]Server)

func GetServer(name string) Server {
	return servers[name]
}

func GetAllServer() map[string]Server {
	return servers
}

func CloseAllServers() {
	for _, server := range servers {
		if server != nil {
			_ = server.Close()
		}
	}
}
