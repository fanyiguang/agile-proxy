package assembly

type Net struct {
	Host     string
	Port     string
	Username string
	Password string
}

func CreateNet(host, port, username, password string) Net {
	return Net{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}
