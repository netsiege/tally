package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
)

// generateNewKey creates a random 32-character API key using readable characters
func generateNewKey() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = charset[int(bytes[i])%len(charset)]
	}

	return string(key), nil
}

// rotateKey generates a new API key and writes it to the platform-specific key file
func rotateKey(task Task) (keyRotationResponse, error) {
	var keyfilePath string

	// Determine key file path based on operating system
	switch runtime.GOOS {
	case "windows":
		keyfilePath = "C:\\Administrator\\.netsiege"
	case "linux":
		keyfilePath = "/root/.netsiege"
	case "darwin":
		keyfilePath = "/Users/akshay/.netsiege"
	default:
		return keyRotationResponse{
			success:        false,
			rotation_error: fmt.Sprintf("unsupported operating system: %s", runtime.GOOS),
		}, nil
	}

	newKey, err := generateNewKey()
	if err != nil {
		return keyRotationResponse{
			success:        false,
			rotation_error: fmt.Sprintf("failed to generate new key: %v", err),
		}, nil
	}

	// Write key with restricted permissions (0600 = read/write for owner only)
	err = os.WriteFile(keyfilePath, []byte(newKey), 0600)
	if err != nil {
		return keyRotationResponse{
			success:        false,
			rotation_error: fmt.Sprintf("failed to write key to file: %v", err),
		}, nil
	}

	return keyRotationResponse{
		success: true,
		new_key: newKey,
	}, nil
}
