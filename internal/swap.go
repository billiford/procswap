package procswap

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/billiford/go-ps"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Swap

// Swap holds functions to implement starting and stopping of batch files.
type Swap interface {
	Path() string
	PID() int
	Start() error
	Kill() error
	Cmd() *exec.Cmd
	ShowOutput(bool)
}

type swap struct {
	cmd        *exec.Cmd
	path       string
	ps         ps.Ps
	showOutput bool
}

// NewSwap returns and implementation of Swap.
func NewSwap(path string) Swap {
	return &swap{
		ps:   ps.New(),
		path: path,
	}
}

// CMD returns the underlying cmd.
func (s *swap) Cmd() *exec.Cmd {
	return s.cmd
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
	// Get the command's stdout pipe so we can show the command's output
	// if requested.
	r, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	// Set the command error to be it's stdout.
	cmd.Stderr = cmd.Stdout
	// Create a scanner which scans r in a line-by-line fashion.
	scanner := bufio.NewScanner(r)
	// Run a thread that is constantly reading the output of a swap until
	// there is no more output :).
	go func(s *swap) {
		// Read line by line and process it.
		for scanner.Scan() {
			line := scanner.Text()
			// Only show the output if the user has requested it.
			if s.showOutput {
				fmt.Println(line)
			}
		}
	}(s)
	// Start the command.
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

// ShowOutput sets the showOutput boolean for this swap.
func (s *swap) ShowOutput(showOutput bool) {
	s.showOutput = showOutput
}
