package utils

import (
	"bytes"
	"os"
	"os/exec"
)

func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	// fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}
