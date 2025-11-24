package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// getTasks retrieves tasks from a file or remote source
func getTasks() (Tasks, error) {
	// TODO: getTasksFromScoreKeeper()
	tasks, err := getTasksFromFile("tasks.json")
	if err != nil {
		fmt.Println("Failed to get tasks:", err)
		return Tasks{}, err
	}
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
