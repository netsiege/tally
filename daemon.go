package main

import (
	"context"
	"fmt"
	"time"
)

// runDaemon contains the core daemon logic with graceful shutdown support
func RunDaemon(ctx context.Context) error {
	fmt.Println("Tally Beacon Service Starting...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Run first iteration immediately
	if err := executeTaskCycle(); err != nil {
		fmt.Printf("Error in task cycle: %v\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown signal received, stopping gracefully...")
			return nil
		case <-ticker.C:
			if err := executeTaskCycle(); err != nil {
				fmt.Printf("Error in task cycle: %v\n", err)
			}
		}
	}
}

// executeTaskCycle performs one iteration of the task processing loop
func executeTaskCycle() error {
	tasks, err := getTasks()
	if err != nil {
		fmt.Println("Error getting tasks:", err)
		return err
	}

	topTask, err := getTopTask(tasks)
	if err != nil {
		fmt.Println("No tasks to execute:", err)
		return err
	}

	controlResp, keyRotResp, err := executeTask(topTask)
	if err != nil {
		fmt.Println("Error executing task:", err)
		return err
	}

	if topTask.TaskType == "check_control" {
		fmt.Printf("Control Check Response: %+v\n", controlResp)
	} else if topTask.TaskType == "rotate_key" {
		fmt.Printf("Key Rotation Response: %+v\n", keyRotResp)
	}

	return nil
}
