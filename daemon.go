package main

import (
	"context"
	"time"
)

// runDaemon contains the core daemon logic with graceful shutdown support
func RunDaemon(ctx context.Context) error {
	LogInfo("Tally Beacon Service Starting...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Run first iteration immediately
	if err := executeTaskCycle(); err != nil {
		LogError("Error in task cycle: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			LogInfo("Shutdown signal received, stopping gracefully...")
			return nil
		case <-ticker.C:
			if err := executeTaskCycle(); err != nil {
				LogError("Error in task cycle: %v", err)
			}
		}
	}
}

// executeTaskCycle performs one iteration of the task processing loop
func executeTaskCycle() error {
	tasks, err := getTasks()
	if err != nil {
		LogError("Error getting tasks: %v", err)
		return err
	}

	topTask, err := getTopTask(tasks)
	if err != nil {
		LogInfo("No tasks to execute: %v", err)
		return err
	}

	// Get the old key before executing the task (important for rotate_key which changes the key)
	oldKey, err := getKey()
	if err != nil {
		LogError("Error getting key: %v", err)
		return err
	}

	controlResp, keyRotResp, err := executeTask(topTask)
	if err != nil {
		LogError("Error executing task: %v", err)
		return err
	}

	if topTask.TaskType == "check_control" {
		err = submitTaskResult(controlResp, oldKey)
		if err != nil {
			LogError("Error submitting check_control response: %v", err)
			return err
		}
	} else if topTask.TaskType == "rotate_key" {
		err = submitKeyRotationResult(keyRotResp, oldKey)
		if err != nil {
			LogError("Error submitting rotate_key response: %v", err)
			return err
		}
	}

	return nil
}
