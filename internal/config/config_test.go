package config

import (
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
