package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	RSS struct {
		UpdatePeriod int      `json:"update_period"`
		Feeds        []string `json:"feeds"`
	} `json:"rss"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
