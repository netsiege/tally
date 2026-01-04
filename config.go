package main

import (
	"fmt"
	"os"
	"runtime"
)

// Configuration variables for the Tally beacon service
var (
	INTERVAL = 60 // in seconds
	ENDPOINT = "http://10.100.7.8:8000"
)

func GetKeyFilePath() (string, error) {
	var keyfilePath string
	switch runtime.GOOS {
	case "windows":
		keyfilePath = "C:\\Users\\Administrator\\.netsiege"
	case "linux":
		keyfilePath = "/root/.netsiege"
	case "darwin":
		keyfilePath = "/Users/akshay/.netsiege"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	_, err := os.Stat(keyfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get key file info: %v", err)
	}

	return keyfilePath, nil
}

func GetEndpointURL(path string) string {
	return "http://" + ENDPOINT + "/api/" + path
}

func getKey() (string, error) {
	keyFilePath, err := GetKeyFilePath()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(keyFilePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
