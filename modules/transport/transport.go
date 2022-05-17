package transport

type Transport interface {
	Transport(ip, port string) (err error)
}
