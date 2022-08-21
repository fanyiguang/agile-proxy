package model

type Base struct {
}

type Net struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Identity struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type PipelineInfos struct {
	PipelineInfo []PipelineInfo
}

type PipelineInfo struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}
