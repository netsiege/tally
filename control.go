package main

import (
	"os"
	"strings"
)

// checkControl verifies the existence and accessibility of a control file
// and returns its contents if successful
func checkControl(task Task) (controlCheckResponse, error) {
	// Handles these cases:
	// 1. file does not exist
	// 2. file exists but cannot be accessed (permission denied)
	// 3. file exists but is too large (over 1MB)
	// 4. file exists but is empty
	// 5. file exists and is accessible, return content

	fileStats, err := os.Stat(task.FilePath)
	if os.IsNotExist(err) {
		return controlCheckResponse{
			Success:     false,
			FilePath:    task.FilePath,
			FileExists:  false,
			AccessError: "[NOTEXIST] - file does not exist",
		}, nil
	}

	if err != nil {
		return controlCheckResponse{
			Success:     false,
			FilePath:    task.FilePath,
			FileExists:  false,
			AccessError: err.Error(),
		}, nil
	}

	// check file size (1 MB limit)
	if fileStats.Size() > 1*1024*1024 {
		return controlCheckResponse{
			Success:     false,
			FilePath:    task.FilePath,
			FileExists:  true,
			AccessError: "[MAXSIZE] - file too large (over 1MB)",
		}, nil
	}

	if fileStats.Size() == 0 {
		return controlCheckResponse{
			Success:     false,
			FilePath:    task.FilePath,
			FileExists:  true,
			AccessError: "[EMPTY] - file is empty",
		}, nil
	}

	content, err := os.ReadFile(task.FilePath)
	if err != nil {
		return controlCheckResponse{
			Success:     false,
			FilePath:    task.FilePath,
			FileExists:  true,
			AccessError: err.Error(),
		}, nil
	}

	// Clean the content by removing newlines
	cleaned := strings.ReplaceAll(string(content), "\n", "")

	return controlCheckResponse{
		Success:     true,
		FilePath:    task.FilePath,
		FileExists:  true,
		FileContent: cleaned,
	}, nil
}
