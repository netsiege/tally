package main

import (
	"crypto/rand"
	"fmt"
	"os"
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
	keyfilePath, err := GetKeyFilePath()
	if err != nil {
		return keyRotationResponse{
			Success:       false,
			RotationError: fmt.Sprintf("failed to get key file path: %v", err),
		}, nil
	}

	newKey, err := generateNewKey()
	if err != nil {
		return keyRotationResponse{
			Success:       false,
			RotationError: fmt.Sprintf("failed to generate new key: %v", err),
		}, nil
	}

	// Write key with restricted permissions (0600 = read/write for owner only)
	_, err = os.Stat(keyfilePath)
	if err != nil {
		return keyRotationResponse{
			Success:       false,
			RotationError: fmt.Sprintf("error checking if key file exists: %v", err),
		}, nil
	}

	if os.IsNotExist(err) {
		// Create the file if it does not exist
		emptyFile, err := os.Create(keyfilePath)
		os.Chmod(keyfilePath, 0600)
		if err != nil {
			return keyRotationResponse{
				Success:       false,
				RotationError: fmt.Sprintf("failed to create key file: %v", err),
			}, nil
		}
		emptyFile.Close()
	}

	err = os.WriteFile(keyfilePath, []byte(newKey), 0600)
	if err != nil {
		return keyRotationResponse{
			Success:       false,
			RotationError: fmt.Sprintf("failed to write key to file: %v", err),
		}, nil
	}

	return keyRotationResponse{
		Success: true,
		NewKey:  newKey,
	}, nil
}
