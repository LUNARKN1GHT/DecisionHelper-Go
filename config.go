package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LLMConfig 用户配置的 LLM 接入信息，存储于本地，不进入代码库
type LLMConfig struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
}

func defaultLLMConfig() LLMConfig {
	return LLMConfig{
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	}
}

func llmConfigFilePath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "llm_config.json"), nil
}

func loadLLMConfig() (LLMConfig, error) {
	path, err := llmConfigFilePath()
	if err != nil {
		return defaultLLMConfig(), nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaultLLMConfig(), nil
	}
	if err != nil {
		return defaultLLMConfig(), nil
	}
	var cfg LLMConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultLLMConfig(), nil
	}
	return cfg, nil
}

func saveLLMConfig(cfg LLMConfig) error {
	path, err := llmConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// 0600: 仅当前用户可读，保护 API Key
	return os.WriteFile(path, data, 0600)
}
