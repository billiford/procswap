package procswap

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/billiford/go-ps"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Swap

// Swap interface holds all fields available to a running swap command.
type Swap interface {
	Cmd() *exec.Cmd
	Path() string
	PID() int
	Stop([]ps.Process) error
	Start() error
}

type swap struct {
	cmd  *exec.Cmd
	path string
	pid  int
}

func NewSwap(path string) Swap {
	return &swap{
		path: path,
	}
}

// PID returns the swaps process id.
func (s *swap) PID() int {
	return s.pid
}

// Path returns the full path of the swap script.
func (s *swap) Path() string {
	return s.path
}

// Cmd returns the underlying exec command for the swap.
func (s *swap) Cmd() *exec.Cmd {
	return s.cmd
}

// Start kicks off a given swap. It just executes the swap process as a command.
func (s *swap) Start() error {
	cmd := exec.Command(s.path)

	err := cmd.Start()
	if err != nil {
		return err
	}

	s.cmd = cmd
	s.pid = cmd.Process.Pid

	return nil
}

// Stop quits a "swap process" or script that is running.
// We need to pass in the currently running processes in order to kill their child processes.
func (s *swap) Stop(processes []ps.Process) error {
	err := killChildProcesses(processes, s.pid)
	if err != nil {
		return fmt.Errorf("error killing child processes for %s: %w", s.cmd.Path, err)
	}

	err = s.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("error killing processes %s: %w", s.cmd.Path, err)
	}

	return nil
}

// killChildProcesses kills all processes that have a parent process ID
// of the process ID passed in.
func killChildProcesses(processes []ps.Process, pid int) error {
	for _, process := range processes {
		if process.PPid() == pid {
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
