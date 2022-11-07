package aria

type Client struct {
	Config
}

func (c *Client) Run(port int) {
	c.RPCListenPort = port

}
