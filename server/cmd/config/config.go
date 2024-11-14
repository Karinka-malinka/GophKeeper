package config

import (
	"encoding/json"
	"os"
)

type ConfigData struct {
	RunAddrgRPS string
	RunAddrREST string
	DatabaseDSN string
	PathKeys    string
}

type ConfigToken struct {
	SecretKeyForToken string `json:"secret_key_for_token"`
	TokenExpiresAt    uint   `json:"token_expires_at"`
	Key               string `json:"key_cripto"`
}

func NewConfig() *ConfigData {
	return &ConfigData{}
}

// GetConfig получает данные конфигурации из файла
func (conf *ConfigData) GetConfig(file string) (*ConfigToken, error) {

	var cfg *ConfigToken

	configToken, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configToken, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
