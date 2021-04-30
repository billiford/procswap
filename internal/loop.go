package procswap

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/billiford/go-ps"
	"github.com/eiannone/keyboard"
	"github.com/logrusorgru/aurora"
)

var (
	defaultPollInterval = 10
	// CurrentSwapOutputIndex holds the index of the swap that is currently
	// printing its output to std out.
	currentSwapOutputIndex = -1
)

// Loop is the interface that runs indefinitely.
type Loop interface {
	Run()
	WithActionsEnabled(bool)
	WithLimit(int)
	WithPollInterval(int)
	WithPriorities([]os.FileInfo)
	WithPriorityScript(string)
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
	// a script that will run when any priority starts
	priorityScript string
	// ps is the interface for listing processes
	ps ps.Ps
	// if the swap scripts have been started or not.
	started bool
	// list of paths to swap scripts.
	swaps []Swap
	// list of currently running swaps.
	runningSwaps []Swap
	// actionsEnabled defines if actions are enabled or not.
	actionsEnabled bool
	// actions is a map of key input to action.
	actions map[rune]action
}

// action holds a key input description and func to call when pressed.
type action struct {
	Description string
	F           func()
}

// NewLoop returns a new Loop.
func NewLoop() Loop {
	// Define the loop.
	loop := &loop{
		swaps:        []Swap{},
		priorities:   []os.FileInfo{},
		limit:        0,
		loopCount:    0,
		ps:           ps.New(),
		pollInterval: defaultPollInterval,
		runningSwaps: []Swap{},
	}
	// Define the actions for the loop. Perhaps this should be defined
	// in main and we should provide a `WithActions(...)` setter function.
	actions := map[rune]action{
		's': {
			Description: "switch console output of swap processes",
			F:           loop.switchOutput,
		},
	}
	// Set the actions for the loop.
	loop.actions = actions

	return loop
}

// WithActionsEnabled enables or disables actions.
func (l *loop) WithActionsEnabled(actionsEnabled bool) {
	l.actionsEnabled = actionsEnabled
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

// WithPriorityScript sets the priority script for the loop.
func (l *loop) WithPriorityScript(priorityScript string) {
	l.priorityScript = priorityScript
}

// WithPs sets the package that will list windows processes.
func (l *loop) WithPs(ps ps.Ps) {
	l.ps = ps
}

// WithSwaps sets the swap scripts/executables for the loop.
func (l *loop) WithSwaps(swaps []Swap) {
	l.swaps = swaps
}

// switchOutput switches the output of running swaps to std out.
func (l *loop) switchOutput() {
	// If there are no currently running swaps, just log this and return.
	if len(l.runningSwaps) == 0 {
		logInfo(fmt.Sprintf("%s no running swaps; ignoring", aurora.Magenta("action")))

		return
	}
	// Hide all outputs.
	for _, swap := range l.runningSwaps {
		swap.ShowOutput(false)
	}
	// Increase the current swap output index.
	currentSwapOutputIndex++
	// If we've gone through all the swap outputs, just hide all outputs.
	if len(l.runningSwaps) <= currentSwapOutputIndex {
		logInfo(fmt.Sprintf("%s hiding all swap output", aurora.Magenta("action")))
		// Reset index to procswap.
		currentSwapOutputIndex = -1

		return
	}
	// Get the output for the current swap output index.
	swap := l.runningSwaps[currentSwapOutputIndex]
	// Let the user know we're showing output for this particular swap.
	logInfo(fmt.Sprintf("%s showing output for %s", aurora.Magenta("action"), aurora.Bold(swap.Path())))

	swap.ShowOutput(true)
}

// Run runs the main loop. It gathers all "priority processes" and runs any swap processes
// when any of these priority processes are not running.
//
// If any priority process starts, all swap processes are killed. When all priority processes
// stop, all swap processes are kicked off again.
func (l *loop) Run() {
	if l.actionsEnabled {
		// Inform the user of the inputs allowed.
		l.printInputDescriptions()
		// Listen for key input in the background.
		go l.listenForKeyInput()
	}

	// Main loop.
	for {
		if l.done() {
			break
		}

		l.run()
		l.wait()
	}
}

func (l *loop) printInputDescriptions() {
	for key, action := range l.actions {
		logInfo(fmt.Sprintf("%s press %s to %s", aurora.Magenta("action"), aurora.BgMagenta(string(key)), action.Description))
	}
}

// listenForKeyInput listens for any key input forever.
// If the input is mapped to some action, procwap will perform this action.
func (l *loop) listenForKeyInput() {
	// Loop forever.
	for {
		// Get key input, for example user has pressed 's'.
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			// Show a warning that there was an error getting key input.
			logWarn(fmt.Sprintf("error getting key input: %s", err.Error()))
			// Continue so this is non-blocking.
			continue
		}
		// If there's an actionable function mapped to this key, run it!
		if _, ok := l.actions[char]; ok {
			l.actions[char].F()
		}
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
		logInfo(fmt.Sprintf("%s %s", aurora.Yellow("priority"), aurora.Bold(strings.Join(runningPriorities, ", "))))

		// It might make sense to set swap scripts to either started or not inside their functions,
		// but I think ths is more explicit.
		l.stop()
		l.stopSwaps()
		l.startPriorityScript()
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
	// Loop through and kill the running swaps.
	for _, swap := range l.runningSwaps {
		logInfo(fmt.Sprintf("%s %s...", aurora.Red("stop"), aurora.Bold(swap.Path())), false)

		err := swap.Kill()
		if err != nil {
			logFailed()
			logError(err.Error())

			continue
		}

		logOK()
	}
	// Since we're shutting down everything, reset the currently running commands.
	l.runningSwaps = []Swap{}
}

// startPriorityScript starts a given priority script. It waits for the command to complete, which
// is different than swaps which are started then stopped if a priority process begins running.
func (l *loop) startPriorityScript() {
	// If no priority script is set, just return.
	if l.priorityScript == "" {
		return
	}

	logInfo(fmt.Sprintf("%s %s...", aurora.Magenta("priority script"), aurora.Bold(l.priorityScript)), false)

	cmd := exec.Command(l.priorityScript)

	err := cmd.Run()
	if err != nil {
		// If there is an error running the priority script, just log it and let the loop continue.
		logFailed()
		logError(err.Error())
	} else {
		logOK()
	}
}
