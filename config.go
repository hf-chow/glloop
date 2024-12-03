package main

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DBURL				string
	CurrentUsername		string
}

const configFilename = ".glloopconfig.json"

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home + "/" + configFilename
	return path, nil
}

func ReadConfig() (Config, error) {
	filePath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	dat, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal(dat, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
