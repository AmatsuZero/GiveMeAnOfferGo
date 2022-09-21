package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProxyConfig struct {
	http, https string
}

type UserConfig struct {
	PathDownloader string       `json:"SaveVideoDir"`
	ConfigProxy    *ProxyConfig `json:"config_proxy"`
	savePath       string
}

func defaultConfig(savePath string) *UserConfig {
	base, _ := os.UserHomeDir()
	return &UserConfig{
		PathDownloader: filepath.Join(base, "Downloads"),
		ConfigProxy:    nil,
		savePath:       savePath,
	}
}

func NewConfig(path string) (*UserConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig(path), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return defaultConfig(path), err
	}
	defer file.Close()
	var tmp UserConfig
	err = json.NewDecoder(file).Decode(&tmp)
	return &tmp, err
}

func (c *UserConfig) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
