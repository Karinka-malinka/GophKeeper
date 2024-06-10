package config

type ConfigData struct {
	RunAddr         string
	ConfigASR       string
	DatabaseDSN     string
	PathFileStorage string
}

func NewConfig() *ConfigData {
	return &ConfigData{}
}
