package normal

type Config struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	ClientNames []string `json:"client_name"`
}
