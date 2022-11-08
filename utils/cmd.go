package utils

import (
	"bytes"
	"errors"
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

func CreateFileIfNotExist(p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		f, e := os.Create(p)
		if e != nil {
			return e
		}
		f.Close()
	}
	return nil
}
