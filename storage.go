package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

func dataDir() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support")
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "DecisionHelper"), nil
}

func dataFilePath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "decisions.json"), nil
}

func loadDecisions() ([]Decision, error) {
	path, err := dataFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Decision{}, nil
	}
	if err != nil {
		return nil, err
	}
	var decisions []Decision
	if err := json.Unmarshal(data, &decisions); err != nil {
		return nil, err
	}
	return decisions, nil
}

func saveDecisions(decisions []Decision) error {
	path, err := dataFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(decisions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
