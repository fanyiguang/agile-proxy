package https

type Config struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	CrtPath  string `json:"crt_path"`
	KeyPath  string `json:"key_path"`
}
