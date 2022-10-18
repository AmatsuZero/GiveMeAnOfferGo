package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Cli struct {
	rootCmd *cobra.Command
}

func NewCli() *Cli {
	cli := &Cli{}
	rootCmd := &cobra.Command{
		Use:   "m3u8-download",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}
	rootCmd.AddCommand(versionCmd)

	cli.rootCmd = rootCmd

	return cli
}

func (c *Cli) Execute() error {
	return c.rootCmd.Execute()
}
