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
			success:      false,
			file_path:    task.FilePath,
			file_exists:  false,
			access_error: "[NOTEXIST] - file does not exist",
		}, nil
	}

	if err != nil {
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  false,
			access_error: err.Error(),
		}, nil
	}

	// check file size (1 MB limit)
	if fileStats.Size() > 1*1024*1024 {
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  true,
			access_error: "[MAXSIZE] - file too large (over 1MB)",
		}, nil
	}

	if fileStats.Size() == 0 {
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  true,
			access_error: "[EMPTY] - file is empty",
		}, nil
	}

	content, err := os.ReadFile(task.FilePath)
	if err != nil {
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  true,
			access_error: err.Error(),
		}, nil
	}

	// Clean the content by removing newlines
	cleaned := strings.ReplaceAll(string(content), "\n", "")

	return controlCheckResponse{
		success:      true,
		file_path:    task.FilePath,
		file_exists:  true,
		file_content: cleaned,
	}, nil
}
