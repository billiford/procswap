package procswap

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-ps"
)

var (
	defaultWaitPeriodSeconds = 10
	runningSwaps             []swap
	started                  bool
	firstLoop                bool
)

type swap struct {
	cmd *exec.Cmd
	pid int
}

// Loop is the interface that runs indefinitely.
type Loop interface {
	Run()
	WithSwaps([]string)
	WithPriorities([]os.FileInfo)
}

// loop holds the priority executables and swap processes defined at startup.
type loop struct {
	swaps      []string
	priorities []os.FileInfo
}

// NewLoop returns a new Loop.
func NewLoop() Loop {
	return &loop{
		swaps:      []string{},
		priorities: []os.FileInfo{},
	}
}

func (l *loop) WithSwaps(swaps []string) {
	l.swaps = swaps
}

func (l *loop) WithPriorities(priorities []os.FileInfo) {
	l.priorities = priorities
}

// Run runs the main loop. It gathers all "priority processes" and runs any swap processes
// when any of these priority processes are not running.
//
// If any priority process starts, all swap processes are killed. When all priority processes
// stop, all swap processes are kicked off again.
func (l *loop) Run() {
	for {
		l.swap()
		wait()

		firstLoop = false
	}
}

func wait() {
	time.Sleep(time.Duration(defaultWaitPeriodSeconds) * time.Second)
}

// swap running swap processes for priority executables or
// start the swaps if no priority process is running and they have not
// already been started.
func (l *loop) swap() {
	// This seems to be a fairly cheap call to check the running processes.
	// It would be nice to just have a watch.
	processes, err := ps.Processes()
	if err != nil {
		logError(fmt.Sprintf("error listing currently running processes: %s", err.Error()))
	}

	// Make a map of the processes so the lookup is O(1).
	processMap := map[string]bool{}
	for _, process := range processes {
		processMap[process.Executable()] = true
	}

	priorities := []string{}
	// Check if an executable has started that we want to take priority over
	// our swap processes.
	for _, priority := range l.priorities {
		if processMap[priority.Name()] {
			priorities = append(priorities, priority.Name())
		}
	}

	// Remove any duplicates and sort the priorities.
	priorities = removeDuplicates(priorities)
	sort.Strings(priorities)

	switch {
	case len(priorities) > 0 && !started && firstLoop:
		// Do this if it is our first loop and priority processes are already running.
		logWarn(fmt.Sprintf("not starting swap processes, priority processes already running: %s",
			aurora.Bold(strings.Join(priorities, ", "))))
	case len(priorities) > 0 && started:
		// Do this if there are any priorities started and we need to stop all running swap processes.
		started = false

		if len(runningSwaps) > 0 {
			logInfo(fmt.Sprintf("priority processes started: %s",
				aurora.Bold(strings.Join(priorities, ", "))))

			l.stopSwaps(processes)
		}
	case len(priorities) == 0 && !started:
		// Do this when there are no priotiies started and we need to start all the swap processes.
		started = true

		logInfo("no priority processes running, starting all swap processes")

		l.startSwaps()
	}
}

func (l *loop) startSwaps() {
	for _, s := range l.swaps {
		// Print this without a newline at the end since we'll be printing the status later.
		logInfo(fmt.Sprintf("%s %s...", aurora.Green("start"), aurora.Bold(s)), false)
		cmd := exec.Command(s)

		err := cmd.Start()
		if err != nil {
			printStatus(err)
			logError(fmt.Sprintf("error starting swap process %s: %s", s, err.Error()))

			continue
		}

		printStatus(err)

		s := swap{
			cmd: cmd,
			pid: cmd.Process.Pid,
		}
		runningSwaps = append(runningSwaps, s)
	}
}

// stopSwaps kills all running swap processes. It finds any child processes
// started by the swap process and attempts to kill those, then kills the main
// process.
//
// We should really build a process ID tree here, but for now the killing of child
// processes is pretty simple.
func (l *loop) stopSwaps(processes []ps.Process) {
	if len(runningSwaps) == 0 {
		logWarn("no swap processes to stop")
		return
	}

	// Store a list of pids that were unsuccessfully killed to add to the list
	// of currently running swap processes.
	pids := map[int]bool{}

	for _, swap := range runningSwaps {
		logInfo(fmt.Sprintf("%s %s...", aurora.Red("stop"), aurora.Bold(swap.cmd.Path)), false)

		err := killChildProcesses(processes, swap.pid)
		if err != nil {
			printStatus(err)
			logError(fmt.Sprintf("error killing child processes for %s: %s", swap.cmd.Path, err.Error()))

			pids[swap.pid] = true

			continue
		}

		err = swap.cmd.Process.Kill()
		if err != nil {
			printStatus(err)
			logError(fmt.Sprintf("error killing processes %s: %s", swap.cmd.Path, err.Error()))

			pids[swap.pid] = true

			continue
		}

		printStatus(err)
	}

	tmpRunningSwaps := []swap{}

	// If any swap processes failed to stop, add them here.
	// TODO we need to figure out a way to come back and retry killing these processes.
	for _, swap := range runningSwaps {
		if pids[swap.pid] {
			tmpRunningSwaps = append(tmpRunningSwaps, swap)
		}
	}

	// Since we're shutting down everything, reset the currently running commands.
	runningSwaps = tmpRunningSwaps
}

// printStatus is a small helper function to either print "OK"
// or "FAILED" in appropriate colors based on an error input.
func printStatus(err error) {
	if err == nil {
		log.Printf(" %s", aurora.Green("OK"))
	} else {
		log.Printf(" %s", aurora.Red("FAILED"))
	}
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

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
