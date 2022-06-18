package client

import jsoniter "github.com/json-iterator/go"

var (
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
	clients = make(map[string]Client)
)

func GetClient(name string) (t Client) {
	return clients[name]
}

func GetAllClients() map[string]Client {
	return clients
}

func CloseAllClients() {
	for _, client := range clients {
		if client != nil {
			_ = client.Close()
		}
	}
}
