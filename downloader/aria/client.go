package aria

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"GiveMeAnOffer/utils"

	"github.com/skratchdot/open-golang/open"
)

type Client struct {
	Config             *Config
	IsRPCServerRunning bool
}

func (c *Client) RunRPCServer() error {
	if c.IsRPCServerRunning {
		return nil
	}

	p, err := c.Config.GenerateConfigFile()
	if err != nil {
		return err
	}

	c.IsRPCServerRunning = true

	_, err = utils.Cmd("aria2c", []string{
		"--conf-path",
		p,
	})

	c.IsRPCServerRunning = false

	return err
}

func (c *Client) RunLocal(port int) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(c.Config.ConfigDir, "static", "index.html")
		http.ServeFile(w, r, p)
	})

	p := fmt.Sprintf(":%v", port)
	open.Run("http://localhost" + p)
	err := http.ListenAndServe(p, nil)
	log.Fatalf("启动本地服务失败: %v", err)

	return nil
}
