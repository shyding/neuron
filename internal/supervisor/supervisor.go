package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// ProcessSpec describes a child process to be managed
type ProcessSpec struct {
	Name       string
	Executable string
	Args       []string
	Dir        string
	Env        []string
}

// Supervisor manages the lifecycle of child processes
type Supervisor struct {
	mu        sync.Mutex
	processes map[string]*exec.Cmd
	stopped   bool
}

// NewSupervisor creates a new process supervisor
func NewSupervisor() *Supervisor {
	return &Supervisor{
		processes: make(map[string]*exec.Cmd),
	}
}

// StartProcess launches and supervises a child process
func (s *Supervisor) StartProcess(spec ProcessSpec) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return fmt.Errorf("supervisor stopped")
	}

	execPath := spec.Executable
	if !filepath.IsAbs(execPath) {
		if abs, err := filepath.Abs(execPath); err == nil {
			execPath = abs
		}
	}

	if _, err := os.Stat(execPath); err != nil {
		fmt.Printf("[Supervisor] Warning: executable %s not found, skipping\n", execPath)
		return nil
	}

	cmd := exec.Command(execPath, spec.Args...)
	if spec.Dir != "" {
		cmd.Dir = spec.Dir
	}
	if len(spec.Env) > 0 {
		cmd.Env = append(os.Environ(), spec.Env...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", spec.Name, err)
	}

	fmt.Printf("[Supervisor] Started %s [PID=%d]\n", spec.Name, cmd.Process.Pid)
	s.processes[spec.Name] = cmd

	go func(name string, c *exec.Cmd) {
		err := c.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.stopped {
			fmt.Printf("[Supervisor] Process %s exited: %v\n", name, err)
		}
	}(spec.Name, cmd)

	return nil
}

// StopAll terminates all managed child processes
func (s *Supervisor) StopAll() {
	s.mu.Lock()
	s.stopped = true
	defer s.mu.Unlock()

	fmt.Println("[Supervisor] Shutting down child processes...")
	for name, cmd := range s.processes {
		if cmd != nil && cmd.Process != nil {
			fmt.Printf("[Supervisor] Stopping %s [PID=%d]...\n", name, cmd.Process.Pid)
			_ = cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	time.Sleep(500 * time.Millisecond)

	for name, cmd := range s.processes {
		if cmd != nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Process.Release()
		}
		delete(s.processes, name)
	}
}
