package main

import (
	"os/exec"
	"sync"
	"syscall"
)

type Proc struct {
	Dir     string
	Command string
	cmd     *exec.Cmd
	mu      sync.Mutex
}

func (proc *Proc) Start() {
	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.cmd != nil && proc.cmd.Process.Signal(syscall.Signal(0)) == nil {
		return
	}

	proc.cmd = exec.Command("sh", "-c", proc.Command)
	proc.cmd.Dir = proc.Dir

	err := proc.cmd.Start()
	if err != nil {
		panic(err)
	}
}

func (proc *Proc) StopGracefully() {
	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.cmd == nil {
		return
	}

	proc.cmd.Process.Signal(syscall.SIGTERM)

	return
}

func (proc *Proc) Wait() error {
	proc.mu.Lock()
	defer proc.mu.Unlock()

	return proc.cmd.Wait()
}

func (proc *Proc) IsRunning() bool {
	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.cmd == nil {
		return false
	}

	process := proc.cmd.Process
	return process.Signal(syscall.Signal(0)) == nil
}

// Add a sigkill for misbehaving processes?
func (proc *Proc) Stop() error {
	proc.StopGracefully()
	if proc.IsRunning() {
		return proc.Wait()
	}

	return nil
}

func (proc *Proc) Restart() (err error) {
	err = proc.Stop()
	proc.Start()
	return
}
