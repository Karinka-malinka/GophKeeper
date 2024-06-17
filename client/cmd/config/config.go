package config

type ConfigData struct {
	ServerAddr string
}

func NewConfig() *ConfigData {
	return &ConfigData{}
}
