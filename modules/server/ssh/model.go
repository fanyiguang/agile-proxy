package ssh

type DirectForward struct {
	DesAddr string
	DesPort uint32

	OriginAddr string
	OriginPort uint32
}

type Config struct {
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	TransportName string `json:"transport_name"`
	RsaPath       string `json:"rsa_path"`
}
