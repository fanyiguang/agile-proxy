package dialer

import jsoniter "github.com/json-iterator/go"

var (
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
	dialers = make(map[string]Dialer)
)

func GetDialer(name string) Dialer {
	return dialers[name]
}

func GetAllDialer() map[string]Dialer {
	return dialers
}
