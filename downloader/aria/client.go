package aria

import (
	"GiveMeAnOffer/utils"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"path/filepath"
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
			log.Fatal(err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(c.Config.ConfigDir, "static", "index.html")
		http.ServeFile(w, r, p)
	})

	p := fmt.Sprintf(":%v", port)
	open.Run("http://localhost" + p)
	log.Fatal(http.ListenAndServe(p, nil))

	return nil
}
