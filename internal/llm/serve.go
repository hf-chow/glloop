package llm

import (
	"fmt"
	"os/exec"
)

func ServeModel() error{
	err := brewStopOllama()
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", "ollama serve")
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func brewStopOllama() error{
	cmd := exec.Command("bash", "-c", "brew services stop ollama" )
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error when stopping: %v", err)
		return err
	}
	return nil
}
