package main

import (
	"fmt"
	"time"
)

// Flow:
// 1. Get tasks from ScoreKeeper or local file
// 2. Pick top task, as they are ordered by priority (queue)
// 3. Execute task based on type
// 4. Send the appropriate response back to ScoreKeeper
// 5. Repeat every X seconds

func main() {
	fmt.Println("Tally Beacon Service Starting...")

	for {
		tasks, err := getTasks()
		if err != nil {
			fmt.Println("Error getting tasks:", err)
			continue
		}

		topTask, err := getTopTask(tasks)
		if err != nil {
			fmt.Println("No tasks to execute:", err)
			continue
		}

		controlResp, keyRotResp, err := executeTask(topTask)
		if err != nil {
			fmt.Println("Error executing task:", err)
			continue
		}

		if topTask.TaskType == "check_control" {
			fmt.Printf("Control Check Response: %+v\n", controlResp)
		} else if topTask.TaskType == "rotate_key" {
			fmt.Printf("Key Rotation Response: %+v\n", keyRotResp)
		}

		// Sleep for a while before checking for new tasks
		// time.Sleep(time.Duration(INTERVAL) * time.Second)
		time.Sleep(time.Duration(10) * time.Second)
	}
}
