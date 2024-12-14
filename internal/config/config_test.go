package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	home := os.Getenv("HOME")
	defer os.Setenv("HOME", home)

	mockHome := "/mock/home"
	err := os.Setenv("HOME", mockHome)
	if err != nil {
		t.Fatalf("unexpected failure setting env variable: %v", err)
	}

	actual, err := getConfigPath() 
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	expectedPath := filepath.Join(mockHome, configFilename)
	if actual != expectedPath {
		t.Errorf("expected %s, but got %s", expectedPath, actual)
	}
}

func TestReadConfig(t *testing.T) {
	home := os.Getenv("HOME")
	defer os.Setenv("HOME", home)

	mockHome, err := os.MkdirTemp("", "mockhome")
	if err != nil {
		t.Fatalf("faled to create mock home dir: %v", err)
	}

	err = os.Setenv("HOME", mockHome)
	if err != nil {
		t.Fatalf("unexpected failure setting env variable: %v", err)
	}

	mockPath := filepath.Join(mockHome, configFilename)

	_, err = ReadConfig()
	if err == nil {
		t.Errorf("expected error when config file does not exist, but got no error")
	}

	mockConfig := Config{
		DBURL:				"postgres://test@localhost:5432/test",
		CurrentUsername: 	"test",
	}

	mockDat, err := json.Marshal(mockConfig)
	if err != nil {
		t.Fatalf("failed to marsahl mock config: %v", err)
	}

	err = os.WriteFile(mockPath, mockDat, 0644)
	if err != nil {
		t.Fatalf("failed to write mock config file: %v", err)
	}

	actualConfig, err := ReadConfig()
	if err != nil {
		t.Errorf("unexpected error when reading actual config: %v", err)
	}

	if actualConfig.DBURL != mockConfig.DBURL || actualConfig.CurrentUsername != mockConfig.CurrentUsername {
		t.Errorf("expected config %+v, but got %+v", mockConfig, actualConfig)
	}

	err = os.WriteFile(mockPath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("failed to wrtie invalid json: %v", err)
	}
	_, err = ReadConfig()
	if err == nil {
		t.Errorf("expected error when config file contains invalid JSON, but get no error instead: %v", err)
	}
}
