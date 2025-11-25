package main

import (
	"context"

	"github.com/kardianos/service"
)

// program implements the service.Interface for cross-platform service management
type program struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Start is called by the service manager when the service starts
// This method must be non-blocking and return quickly
func (p *program) Start(s service.Service) error {
	// Create a cancellable context for graceful shutdown
	p.ctx, p.cancel = context.WithCancel(context.Background())

	// Start the daemon in a goroutine so Start() returns immediately
	go func() {
		RunDaemon(p.ctx)
	}()

	return nil
}

// Stop is called by the service manager when the service stops
// This triggers graceful shutdown of the daemon
func (p *program) Stop(s service.Service) error {
	// Cancel the context, which triggers ctx.Done() in RunDaemon
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}
