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
	Config *Config
}

func (c *Client) RunLocal(port int) error {
	go func() {
		p, err := c.Config.GenerateConfigFile()
		if err != nil {
			log.Fatal(err)
		}

		_, err = utils.Cmd("aria2c", []string{
			"--conf-path",
			p,
		})

		if err != nil {
			log.Fatalf("启动 aria2 失败： %v", err)
		}
	}()

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
