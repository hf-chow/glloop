package main

import (
	"fmt"
	"os/exec"
)

const PORT = 11434
var ENDPOINT string = fmt.Sprintf("http://localhost:%d/api/", int(PORT))

func serveModel() error{
	cmd := exec.Command("bash", "-c",  "ollama serve")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
