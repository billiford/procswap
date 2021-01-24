package procswap

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-ps"
	"github.com/shiena/ansicolor"
)

var (
	defaultPollInterval = 10
)

type swap struct {
	cmd *exec.Cmd
	pid int
}

// Loop is the interface that runs indefinitely.
type Loop interface {
	Run()
	WithLimit(int)
	WithPollInterval(int)
	WithPriorities([]os.FileInfo)
	WithSwaps([]string)
}

// loop holds the priority executables and swap processes defined at startup.
type loop struct {
	swaps        []string
	priorities   []os.FileInfo
	limit        int
	pollInterval int
	loopCount    int
	started      bool
	runningSwaps []swap
}

// NewLoop returns a new Loop.
func NewLoop() Loop {
	return &loop{
		swaps:        []string{},
		priorities:   []os.FileInfo{},
		limit:        0,
		loopCount:    0,
		pollInterval: defaultPollInterval,
		runningSwaps: []swap{},
	}
}

// WithLimit sets a limit on the loop.
func (l *loop) WithLimit(limit int) {
	l.limit = limit
}

// WithPollInterval sets the poll interval on the loop.
func (l *loop) WithPollInterval(pollInterval int) {
	l.pollInterval = pollInterval
}

// WithPriorities sets the priority processes for the loop.
func (l *loop) WithPriorities(priorities []os.FileInfo) {
	l.priorities = priorities
}

// WithSwaps sets the swap scripts/executables for the loop.
func (l *loop) WithSwaps(swaps []string) {
	l.swaps = swaps
}

// Run runs the main loop. It gathers all "priority processes" and runs any swap processes
// when any of these priority processes are not running.
//
// If any priority process starts, all swap processes are killed. When all priority processes
// stop, all swap processes are kicked off again.
func (l *loop) Run() {
	for {
		if l.done() {
			break
		}

		l.loop()
		l.wait()
	}
}

// loop runs the main loop. Swap running processes for priority executables or
// start the swap processes if no priority process is running and they have not
// already been started.
func (l *loop) loop() {
	// This seems to be a fairly cheap call to check the running processes.
	// It would be nice to just have a watch.
	processes, err := ps.Processes()
	if err != nil {
		logError(fmt.Sprintf("error listing currently running processes: %s", err.Error()))

		return
	}

	// List running priorities from the current processes running.
	runningPriorities := l.listRunningPriorities(processes)

	switch {
	case len(runningPriorities) > 0 && !l.started && l.loopCount == 0:
		// It is our first loop and priority processes are already running so log this.
		logWarn(fmt.Sprintf("not starting swap processes, priority processes already running: %s",
			aurora.Bold(strings.Join(runningPriorities, ", "))))
	case len(runningPriorities) > 0 && l.started:
		// Do this if there are any priorities started and we need to stop all running swap processes.
		logInfo(fmt.Sprintf("%s %s", aurora.Yellow("start"), aurora.Bold(strings.Join(runningPriorities, ", "))))

		l.stop()
		l.stopSwaps(processes)
	case len(runningPriorities) == 0 && !l.started:
		// Do this when there are no priorities started and we need to start all the swap processes.
		l.start()
		l.startSwaps()
	}

	l.loopCount++
}

func (l *loop) listRunningPriorities(processes []ps.Process) []string {
	// Make a map of the processes so the lookup is O(1).
	processMap := map[string]bool{}
	for _, process := range processes {
		processMap[process.Executable()] = true
	}

	prioritiesMap := map[string]bool{}
	// Check if an executable has started that we want to take priority over
	// our swap processes.
	for _, priority := range l.priorities {
		if processMap[priority.Name()] {
			prioritiesMap[priority.Name()] = true
		}
	}

	priorities := make([]string, 0, len(prioritiesMap))
	for k := range prioritiesMap {
		priorities = append(priorities, k)
	}

	sort.Strings(priorities)

	return priorities
}

func (l *loop) wait() {
	time.Sleep(time.Duration(l.pollInterval) * time.Second)
}

func (l *loop) stop() {
	l.started = false
}

func (l *loop) start() {
	l.started = true
}

func (l *loop) done() bool {
	if l.limit < 1 {
		return false
	}

	return l.limit == l.loopCount
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
		l.runningSwaps = append(l.runningSwaps, s)
	}
}

// stopSwaps kills all running swap processes. It finds any child processes
// started by the swap process and attempts to kill those, then kills the main
// process.
//
// We should really build a process ID tree here, but for now the killing of child
// processes is pretty simple.
func (l *loop) stopSwaps(processes []ps.Process) {
	if len(l.runningSwaps) == 0 {
		logWarn("no swap processes to stop")

		return
	}

	// Store a list of pids that were unsuccessfully killed to add to the list
	// of currently running swap processes.
	pids := map[int]bool{}

	for _, swap := range l.runningSwaps {
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
	for _, swap := range l.runningSwaps {
		if pids[swap.pid] {
			tmpRunningSwaps = append(tmpRunningSwaps, swap)
		}
	}

	// Since we're shutting down everything, reset the currently running commands.
	l.runningSwaps = tmpRunningSwaps
}

// printStatus is a small helper function to either print "OK"
// or "FAILED" in appropriate colors based on an error input.
func printStatus(err error) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	if err == nil {
		fmt.Fprintf(w, " %s\n", aurora.Green("OK"))
	} else {
		fmt.Fprintf(w, " %s\n", aurora.Red("FAILED"))
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
