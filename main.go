package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kardianos/service"
)

const Version = "1.0.0"
const BuildDate = "2025-11-24"

// Flow:
// 1. Get tasks from ScoreKeeper or local file
// 2. Pick top task, as they are ordered by priority (queue)
// 3. Execute task based on type
// 4. Send the appropriate response back to ScoreKeeper
// 5. Repeat every X seconds

func main() {
	// Service configuration
	svcConfig := &service.Config{
		Name:        "tally",
		DisplayName: "tally Beacon Service",
		Description: "Monitors and executes control tasks from scorekeeper - used for scoring netsiege",
		// WorkingDirectory: "/Users/akshay/Documents/GitHub/tally/tally",
		Arguments: []string{},
	}

	// Create program instance
	prg := &program{}

	// Create service wrapper
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		os.Exit(1)
	}

	// Handle service control commands (install, uninstall, start, stop)
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		switch cmd {
		case "status":
			showStatus()
			return

		case "logs":
			showLogs()
			return

		case "version":
			fmt.Printf("Tally Beacon Service v%s\n", Version)
			fmt.Printf("Build: %s\n", BuildDate)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return

		default:
			err = service.Control(s, cmd)
			if err != nil {
				fmt.Printf("Valid commands: install, uninstall, start, stop, restart, status, logs, version\n")
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	// Run the service (works in both console and service mode)
	err = s.Run()
	if err != nil {
		fmt.Printf("Error running service: %v\n", err)
		os.Exit(1)
	}
}
