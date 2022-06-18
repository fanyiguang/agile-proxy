package transport

import jsoniter "github.com/json-iterator/go"

var (
	transports = make(map[string]Transport)
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
)

func GetTransport(name string) (t Transport) {
	return transports[name]
}

func GetAllTransport() map[string]Transport {
	return transports
}

func CloseAllTransports() {
	for _, transport := range transports {
		if transport != nil {
			_ = transport.Close()
		}
	}
}
