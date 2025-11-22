package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var ENDPOINT string
var HOSTNAME string
var PORT int

// Flow:
/*
1. Get tasks from ScoreKeeper or local file
2. Pick top task, as they are ordered by priority (queue)
3. Execute task based on type
4. Send the appropriate response back to ScoreKeeper
5. Repeat every X seconds
*/

func main() {
	fmt.Println("Tally Beacon Service Starting...")
	tasks, err := getTasks()
	if err != nil {
		fmt.Println("Error getting tasks:", err)
		return
	}

	topTask, err := getTopTask(tasks)
	if err != nil {
		fmt.Println("Error getting top task:", err)
		return
	}

	controlResp, keyRotResp, err := executeTask(topTask)
	if err != nil {
		fmt.Println("Error executing task:", err)
		return
	}

	fmt.Println("Control Check Response:", controlResp)
	fmt.Println("Key Rotation Response:", keyRotResp)
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	TaskType string `json:"type"`
	FilePath string `json:"file_path,omitempty"`
}

func getTasksFromFile(filePath string) (Tasks, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return Tasks{}, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Tasks{}, err
	}

	var tasks Tasks
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return Tasks{}, err
	}

	return tasks, nil

}

func getTasks() (Tasks, error) {
	// getTasksFromScoreKeeper()
	tasks, err := getTasksFromFile("tasks.json")
	if err != nil {
		fmt.Println("Failed to get tasks:", err)
		return Tasks{}, err
	}
	return tasks, nil
}

func getTopTask(tasks Tasks) (Task, error) {
	if len(tasks.Tasks) == 0 {
		return Task{}, fmt.Errorf("no tasks found, we are all caught up!")
	}
	return tasks.Tasks[0], nil
}

func executeTask(task Task) (controlCheckResponse, keyRotationResponse, error) {
	switch task.TaskType {
	case "check_control":
		fmt.Println("Executing check_control task")
		resp, err := checkControl(task)
		if err != nil {
			return controlCheckResponse{}, keyRotationResponse{}, err
		}
		return resp, keyRotationResponse{}, nil

	case "process_file":
		fmt.Println("Executing process_file task on file:", task.FilePath)

	case "rotate_key":
		fmt.Println("Executing rotate_key task")
		resp, err := rotateKey(task)
		if err != nil {
			return controlCheckResponse{}, keyRotationResponse{}, err
		}
		return controlCheckResponse{}, resp, nil
	}

	return controlCheckResponse{}, keyRotationResponse{}, fmt.Errorf("unknown task type: %s", task.TaskType)
}

type keyRotationResponse struct {
	success        bool
	new_key        string
	rotation_error string
}

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

func rotateKey(task Task) (keyRotationResponse, error) {
	var keyfilePath string
	if runtime.GOOS == "windows" {
		keyfilePath = "C:\\Administrator\\.netsiege"
	} else if runtime.GOOS == "linux" {
		keyfilePath = "/root/.netsiege"
	} else if runtime.GOOS == "darwin" {
		keyfilePath = "/Users/akshay/.netsiege"
	}

	newKey, err := generateNewKey()
	if err != nil {
		return keyRotationResponse{
			success:        false,
			rotation_error: fmt.Sprintf("failed to generate new key: %v", err),
		}, nil
	}

	err = os.WriteFile(keyfilePath, []byte(newKey), 0600) // 0600 = read/write for owner only
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

type controlCheckResponse struct {
	success      bool
	file_path    string
	file_exists  bool
	file_content string
	access_error string
}

func checkControl(task Task) (controlCheckResponse, error) {
	// cases:
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
			access_error: "file does not exist",
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

	// check file size
	if fileStats.Size() > 1*1024*1024 { // 1 MB limit
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  true,
			access_error: "file too large (over 1MB)",
		}, nil
	}

	if fileStats.Size() == 0 {
		return controlCheckResponse{
			success:      false,
			file_path:    task.FilePath,
			file_exists:  true,
			access_error: "file is empty",
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

	cleaned := strings.ReplaceAll(string(content), "\n", "")

	return controlCheckResponse{
		success:      true,
		file_path:    task.FilePath,
		file_exists:  true,
		file_content: cleaned,
	}, nil
}
