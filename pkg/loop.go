package procswap

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/billiford/go-ps"
	"github.com/logrusorgru/aurora"
)

var (
	defaultPollInterval = 10
)

// Loop is the interface that runs indefinitely.
type Loop interface {
	Run()
	WithLimit(int)
	WithPollInterval(int)
	WithPriorities([]os.FileInfo)
	WithPs(ps.Ps)
	WithSwaps([]Swap)
}

// loop holds the priority executables and swap processes defined at startup.
type loop struct {
	// limit is the limit the loop (polling windows ps) will run; less than 1 is infinite times.
	limit int
	// internal storage of how many times we've looped.
	loopCount int
	// poll interval sets how much time in seconds we wait before polling the windows processes.
	pollInterval int
	// list of priorities defined at startup.
	priorities []os.FileInfo
	// ps is the interface for listing processes
	ps ps.Ps
	// if the swap scripts have been started or not.
	started bool
	// list of paths to swap scripts.
	swaps []Swap
	// list of currently running swaps.
	runningSwaps []Swap
}

// NewLoop returns a new Loop.
func NewLoop() Loop {
	return &loop{
		swaps:        []Swap{},
		priorities:   []os.FileInfo{},
		limit:        0,
		loopCount:    0,
		ps:           ps.New(),
		pollInterval: defaultPollInterval,
		runningSwaps: []Swap{},
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

// WithPs sets the package that will list windows processes.
func (l *loop) WithPs(ps ps.Ps) {
	l.ps = ps
}

// WithSwaps sets the swap scripts/executables for the loop.
func (l *loop) WithSwaps(swaps []Swap) {
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

		l.run()
		l.wait()
	}
}

// run runs the main loop. Swap running processes for priority executables or
// start the swap processes if no priority process is running and they have not
// already been started.
func (l *loop) run() {
	defer l.incCount()

	// List running priorities from the current processes running.
	runningPriorities := l.listRunningPriorities()

	switch {
	case len(runningPriorities) > 0 && !l.started && l.loopCount == 0:
		// It is our first loop and priority processes are already running so log this.
		logWarn(fmt.Sprintf("not starting swap processes, priority processes already running: %s",
			aurora.Bold(strings.Join(runningPriorities, ", "))))
	case len(runningPriorities) > 0 && l.started:
		// Do this if there are any priorities started and we need to stop all running swap processes.
		logInfo(fmt.Sprintf("%s %s", aurora.Yellow("start"), aurora.Bold(strings.Join(runningPriorities, ", "))))

		// It might make sense to set swap scripts to either started or not inside their functions,
		// but I think ths is more explicit.
		l.stop()
		l.stopSwaps()
	case len(runningPriorities) == 0 && !l.started:
		// Do this when there are no priorities started and we need to start all the swap processes.
		l.start()
		l.startSwaps()
	}
}

func (l *loop) incCount() {
	l.loopCount++
}

// listRunningPriorities takes in a list of currently running
// priorities and makes a list of any user-defined priorities that are
// running.
func (l *loop) listRunningPriorities() []string {
	// This seems to be a fairly cheap call to check the running processes.
	// It would be nice to just have a watch.
	processes, err := l.ps.Processes()
	if err != nil {
		logError(fmt.Sprintf("error listing currently running processes: %s", err.Error()))

		return nil
	}

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

	// Generate a slice of currently running priorities.
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
		logInfo(fmt.Sprintf("%s %s...", aurora.Green("start"), aurora.Bold(s.Path())), false)

		err := s.Start()
		if err != nil {
			logFailed()
			logError(fmt.Sprintf("error starting swap process %s: %s", s.Path(), err.Error()))

			continue
		}

		logOK()

		l.runningSwaps = append(l.runningSwaps, s)
	}
}

// stopSwaps kills all running swap processes. It finds any child processes
// started by the swap process and attempts to kill those, then kills the main
// process.
//
// We should really build a process ID tree here, but for now the killing of child
// processes is pretty simple.
func (l *loop) stopSwaps() {
	// Store a list of pids that were unsuccessfully killed to add to the list
	// of currently running swap processes.
	pids := map[int]bool{}

	for _, swap := range l.runningSwaps {
		logInfo(fmt.Sprintf("%s %s...", aurora.Red("stop"), aurora.Bold(swap.Path())), false)

		err := swap.Kill()
		if err != nil {
			logFailed()
			logError(err.Error())

			continue
		}

		logOK()

		pids[swap.PID()] = true
	}

	tmpRunningSwaps := []Swap{}

	// If any swap processes failed to stop, add them here.
	// TODO we need to figure out a way to come back and retry killing these processes.
	for _, swap := range l.runningSwaps {
		if pids[swap.PID()] {
			tmpRunningSwaps = append(tmpRunningSwaps, swap)
		}
	}

	// Since we're shutting down everything, reset the currently running commands.
	l.runningSwaps = tmpRunningSwaps
}
