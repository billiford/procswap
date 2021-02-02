package procswap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/billiford/go-ps"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Swap

type Swap interface {
	Path() string
	PID() int
	Start() error
	Kill() error
}

type swap struct {
	cmd  *exec.Cmd
	path string
	ps   ps.Ps
}

func NewSwap(path string) Swap {
	return &swap{
		ps:   ps.New(),
		path: path,
	}
}

// Path returns the path of the command.
func (s *swap) Path() string {
	return s.path
}

// PID returns the process ID of the underlying command. If there is no
// command it returns -1.
func (s *swap) PID() int {
	if s.cmd == nil {
		return -1
	}

	return s.cmd.Process.Pid
}

// Start starts a given command.
func (s *swap) Start() error {
	cmd := exec.Command(s.path)
	s.cmd = cmd

	return cmd.Start()
}

// Kill kills all direct child processes of the PID passed in,
// then attempts to kill the PID itself.
func (s *swap) Kill() error {
	if s.cmd == nil {
		return errors.New("no command to kill")
	}

	// Kill all child processes.
	err := s.killChildProcesses()
	if err != nil {
		return fmt.Errorf("error killing child processes for %s: %w", s.path, err)
	}

	// Kill the process.
	err = s.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("error killing processes %s: %w", s.path, err)
	}

	_, err = s.cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("error waiting on process to be killed %s: %w", s.path, err)
	}

	return nil
}

// killChildProcesses kills all processes that have a parent process ID
// of the process ID passed in.
func (s *swap) killChildProcesses() error {
	// List all currently running processes.
	processes, err := s.ps.Processes()
	if err != nil {
		return err
	}

	for _, process := range processes {
		if process.PPid() == s.cmd.Process.Pid {
			p, err := os.FindProcess(process.Pid())
			if err != nil {
				return fmt.Errorf("error finding process %s: %w", process.Executable(), err)
			}

			err = p.Kill()
			if err != nil {
				return fmt.Errorf("error killing process %s: %w", process.Executable(), err)
			}
		}
	}

	return nil
}
