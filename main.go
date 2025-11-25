package main

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
)

// Flow:
// 1. Get tasks from ScoreKeeper or local file
// 2. Pick top task, as they are ordered by priority (queue)
// 3. Execute task based on type
// 4. Send the appropriate response back to ScoreKeeper
// 5. Repeat every X seconds

func main() {
	// Service configuration
	svcConfig := &service.Config{
		Name:             "tally",
		DisplayName:      "tally Beacon Service",
		Description:      "Monitors and executes control tasks from scorekeeper - used for scoring netsiege",
		// WorkingDirectory: "/Users/akshay/Documents/GitHub/tally/tally",
		Arguments:        []string{},
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
		err = service.Control(s, os.Args[1])
		if err != nil {
			fmt.Printf("Valid commands: install, uninstall, start, stop, restart\n")
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
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
