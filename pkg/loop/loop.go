package loop

import (
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

var (
	defaultWaitPeriodSeconds = 1
	runningCmds              []*exec.Cmd
	started                  bool
)

// Loop is the interface that runs.
type Loop interface {
	Run()
	WithBackgroundScripts([]string)
	WithPriorityExecutables([]os.FileInfo)
}

// loop holds the priority executables and background scripts defined at startup.
type loop struct {
	backgroundScripts   []string
	priorityExecutables []os.FileInfo
}

// New returns a new Loop.
func New() Loop {
	return &loop{
		backgroundScripts:   []string{},
		priorityExecutables: []os.FileInfo{},
	}
}

func (l *loop) WithBackgroundScripts(backgroundScripts []string) {
	l.backgroundScripts = backgroundScripts
}

func (l *loop) WithPriorityExecutables(priorityExecutables []os.FileInfo) {
	l.priorityExecutables = priorityExecutables
}

// Run runs the main loop. It gathers all "priority processes" and runs any background scripts
// when any of these priority processes are not running.
func (l *loop) Run() {
	for {
		l.swap()
		wait()
	}
}

func wait() {
	time.Sleep(time.Duration(defaultWaitPeriodSeconds) * time.Second)
}

// swap running background scripts for priority executables or
// start the scripts if no priority process is running and they have not
// already been started.
func (l *loop) swap() {
	// This seems to be a fairly cheap call to check the running processes.
	// It would be nice to just have a watch.
	processes, err := ps.Processes()
	if err != nil {
		// TODO do some fancy logging
		panic(err)
	}

	// Make a map of the processes so the lookup is O(1).
	currentProcessMap := map[string]bool{}
	for _, process := range processes {
		currentProcessMap[process.Executable()] = true
	}

	currentlyRunningPriorityProcesses := []string{}
	// Check if an executable has started that we want to take priority over
	// our scripts.
	for _, priorityExecutable := range l.priorityExecutables {
		if currentProcessMap[priorityExecutable.Name()] {
			currentlyRunningPriorityProcesses = append(currentlyRunningPriorityProcesses, priorityExecutable.Name())
		}
	}

	if len(currentlyRunningPriorityProcesses) > 0 && started {
		started = false

		l.shutdownScripts(currentlyRunningPriorityProcesses, processes)
	} else if len(currentlyRunningPriorityProcesses) == 0 && !started {
		started = true

		l.runScripts()
	}
}

func (l *loop) runScripts() {
	for _, s := range l.backgroundScripts {
		log.Println("starting script", s)
		cmd := exec.Command(s)

		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		runningCmds = append(runningCmds, cmd)
	}
}

func (l *loop) shutdownScripts(priorityExecutables []string,
	currentlyRunningProcesses []ps.Process) {
	priorityExecutables = removeDuplicates(priorityExecutables)
	sort.Strings(priorityExecutables)

	if len(runningCmds) > 0 {
		// TODO add fancy logging.
		log.Printf("The following priority executables were found, so we're shutting down %d scripts: %s\n",
			len(runningCmds), strings.Join(priorityExecutables, ", "))
	}

	for _, cmd := range runningCmds {
		// TODO add fancy logging.
		log.Println("shutting down", cmd.Path)

		err := killChildProcesses(currentlyRunningProcesses, cmd.Process.Pid)
		if err != nil {
			// TODO don't panic, log!
			panic(err)
		}

		err = cmd.Process.Kill()
		if err != nil {
			// TODO don't panic, log!
			panic(err)
		}
	}

	// Since we're shutting down everything, reset the currently running commands.
	runningCmds = []*exec.Cmd{}
}

// killChildProcesses kills all processes that have a parent process ID
// of the process ID passed in.
func killChildProcesses(processes []ps.Process, pid int) error {
	for _, process := range processes {
		if process.PPid() == pid {
			p, err := os.FindProcess(process.Pid())
			if err != nil {
				return err
			}

			err = p.Kill()
			if err != nil {
				return err
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
