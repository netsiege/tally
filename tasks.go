package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// getTasks retrieves tasks from a file or remote source
func getTasks() (Tasks, error) {
	tasks, err := getTasksFromScoreKeeper()
	if err != nil {
		fmt.Println("Failed to get tasks:", err)
		return Tasks{}, err
	}
	// tasks, err := getTasksFromFile("tasks.json")
	// if err != nil {
	// 	fmt.Println("Failed to get tasks:", err)
	// 	return Tasks{}, err
	// }
	return tasks, nil
}

// getTasksFromFile reads and parses tasks from a JSON file
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

// getTasksFromScoreKeeper retrieves tasks from the scorekeeper API
func getTasksFromScoreKeeper() (Tasks, error) {
	tasksEndpoint := GetEndpointURL("/api/tasks")

	key, err := getKey()
	if err != nil {
		return Tasks{}, fmt.Errorf("failed to get authentication key: %v", err)
	}

	req, err := http.NewRequest("GET", tasksEndpoint, nil)
	if err != nil {
		return Tasks{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("User-Agent", "Tally-Beacon/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Tasks{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Tasks{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Tasks{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var tasks Tasks
	err = json.Unmarshal(responseData, &tasks)
	if err != nil {
		return Tasks{}, fmt.Errorf("failed to unmarshal tasks: %v", err)
	}

	return tasks, nil
}

// getTopTask returns the highest priority task (first in queue)
func getTopTask(tasks Tasks) (Task, error) {
	if len(tasks.Tasks) == 0 {
		return Task{}, fmt.Errorf("no tasks found, we are all caught up!")
	}
	return tasks.Tasks[0], nil
}

// executeTask dispatches the task to the appropriate handler based on type
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
		// TODO: Implement process_file handler

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

func submitTaskResult(checkResponse controlCheckResponse, key string) error {
	taskSubmissionEndpoint := GetEndpointURL("/api/claim")

	jsonControlCheckResponse, err := json.Marshal(checkResponse)
	if err != nil {
		LogError("Failure to marshal control check response: %v", err)
		return err
	}

	_, err = AuthenticatedPostRequestWithPayload(taskSubmissionEndpoint, jsonControlCheckResponse, key)
	if err != nil {
		LogError("Failure to submit task result: %v", err)
		return err
	}

	LogInfo("Successfully submitted check_control response")
	return nil
}

func submitKeyRotationResult(keyRotResponse keyRotationResponse, key string) error {
	taskSubmissionEndpoint := GetEndpointURL("/api/rotate_key")

	jsonKeyRotationResponse, err := json.Marshal(keyRotResponse)
	if err != nil {
		LogError("Failure to marshal key rotation response: %v", err)
		return err
	}

	_, err = AuthenticatedPostRequestWithPayload(taskSubmissionEndpoint, jsonKeyRotationResponse, key)
	if err != nil {
		LogError("Failure to submit key rotation result: %v", err)
		return err
	}

	LogInfo("Successfully submitted rotate_key response")
	return nil
}