package main

import (
	"bufio"
	"os"
)

func readFileByLine(filePath string, block func(string, *bool)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	shouldStop := false
	for line, isPrefix, err := r.ReadLine(); err == nil && !isPrefix; line, isPrefix, err = r.ReadLine() {
		s := string(line)
		block(s, &shouldStop)
		if shouldStop {
			break
		}
	}
	return nil
}
