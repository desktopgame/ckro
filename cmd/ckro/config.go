package main

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	ApiKey       string `yaml:"ApiKey"`
	BaseUrl      string `yaml:"BaseUrl"`
	Model        string `yaml:"Model"`
	SystemPrompt string `yaml:"SystemPrompt"`
}

func (c *Config) Reset() {
	c.ApiKey = "lmstudio"
	c.BaseUrl = "http://localhost:1234/v1"
	c.Model = "openai/gpt-oss-20b"
	c.SystemPrompt = "You are friendly assistant."
}

func (c *Config) Update(vaultManager *VaultManager) {
	c.Reset()

	file := filepath.Join(vaultManager.GetRoot(), "ckro.yaml")
	data, err := os.ReadFile(file)
	if err == nil {
		yaml.Unmarshal(data, c)
	}
}
