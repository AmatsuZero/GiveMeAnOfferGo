package app

import (
	"GiveMeAnOffer/downloader/aria"
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
	ConCurrentCnt  int          `json:"ConCurrentCnt"`
	savePath       string
	AriaConfig     *aria.Config
}

func defaultConfig(savePath string) *UserConfig {
	base, _ := os.UserHomeDir()
	config := &UserConfig{
		PathDownloader: filepath.Join(base, "Downloads"),
		ConfigProxy:    nil,
		savePath:       savePath,
		ConCurrentCnt:  5,
	}

	config.AriaConfig = aria.DefaultConfig(filepath.Join(appFolder, "aria2"))
	config.AriaConfig.MaxConCurrentDownload = config.ConCurrentCnt
	config.AriaConfig.SaveDir = config.PathDownloader

	return config
}

func NewConfig(path string) (*UserConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig(path), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return defaultConfig(path), err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			SharedApp.LogError(err.Error())
		}
	}(file)
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
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			SharedApp.LogError(err.Error())
		}
	}(f)
	_, err = f.Write(data)
	return err
}
