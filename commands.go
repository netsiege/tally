package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// showStatus displays the current service status
func showStatus() {
	fmt.Println("Service: Tally")

	switch runtime.GOOS {
	case "darwin":
		// Check launchctl
		cmd := exec.Command("launchctl", "list", "tally")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Status: Not running")
			return
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) > 0 {
			fields := strings.Fields(lines[0])
			if len(fields) >= 3 {
				fmt.Printf("Status: Running\n")
				fmt.Printf("PID: %s\n", fields[0])
			}
		}

	case "linux":
		// Check systemctl
		cmd := exec.Command("systemctl", "is-active", "tally")
		output, err := cmd.CombinedOutput()
		status := strings.TrimSpace(string(output))

		if err != nil || status != "active" {
			fmt.Println("Status: Not running")
		} else {
			fmt.Println("Status: Running")

			// Get PID
			cmd = exec.Command("systemctl", "show", "tally", "--property=MainPID")
			output, _ = cmd.CombinedOutput()
			fmt.Print(string(output))
		}

	case "windows":
		// Check Windows service
		cmd := exec.Command("sc", "query", "tally")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Status: Not running")
		} else {
			if strings.Contains(string(output), "RUNNING") {
				fmt.Println("Status: Running")
			} else {
				fmt.Println("Status: Stopped")
			}
		}
	}

	fmt.Printf("Log file: %s\n", getLogFilePath())
}

// showLogs displays recent log entries
func showLogs() {
	logPath := getLogFilePath()

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Printf("Log file not found: %s\n", logPath)
		return
	}

	// Show last 50 lines
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "Get-Content", logPath, "-Tail", "50")
	} else {
		cmd = exec.Command("tail", "-n", "50", logPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error reading logs: %v\n", err)
	}
}
