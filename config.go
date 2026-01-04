package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

// Configuration variables for the Tally beacon service
var (
	INTERVAL int
	ENDPOINT string
)

// LoadConfig reads configuration from environment variables
// (which are loaded by systemd from EnvironmentFile or set manually)
func LoadConfig() error {
	var err error

	// Read ENDPOINT from environment (required)
	ENDPOINT = os.Getenv("ENDPOINT")
	if ENDPOINT == "" {
		return fmt.Errorf("ENDPOINT environment variable not set (check systemd EnvironmentFile)")
	}

	// Read INTERVAL from environment (optional, defaults to 10)
	intervalStr := os.Getenv("INTERVAL")
	if intervalStr != "" {
		INTERVAL, err = strconv.Atoi(intervalStr)
		if err != nil {
			return fmt.Errorf("invalid INTERVAL value '%s': %v", intervalStr, err)
		}
	} else {
		INTERVAL = 10 // Default to 10 seconds
	}

	return nil
}

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
