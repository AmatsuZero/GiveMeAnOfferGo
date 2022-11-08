package aria

import (
	"GiveMeAnOffer/utils"
	"fmt"
)

type Client struct {
	Config *Config
}

func (c *Client) RunLocal(port int) error {
	c.Config.RPCListenPort = port
	p, err := c.Config.GenerateConfigFile()
	if err != nil {
		return err
	}
	fmt.Println(p)

	_, err = utils.Cmd("aria2c", []string{
		"--conf-path",
		p,
	})

	if err != nil {
		return err
	}

	return err
}
