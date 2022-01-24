package procswap

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

const (
	appName                   = "procswap"
	appUsage                  = "run processes when any prioritized process is not running"
	appUsageText              = "procswap.exe -p <PATH_TO_DIR_FOR_PRIORITIES> -s <PATH_TO_EXECUTABLE>"
	authorName                = "billiford"
	flagDiableActionsName     = "disable-actions"
	flagDiableActionsUsage    = "disable actions (keyboard inputs)"
	flagIgnoreAliases         = "i"
	flagIgnoreName            = "ignore"
	flagIgnoreUsage           = "ignore a priority (case insensitive)"
	flagLimitAliases          = "l"
	flagLimitName             = "limit"
	flagLimitUsage            = "a limit to a number of times the loop runs (0 = infinite)"
	flagLimitValue            = 0
	flagPollIntervalAliases   = "pi"
	flagPollIntervalName      = "poll-interval"
	flagPollIntervalUsage     = "time in seconds to wait to poll for running processes"
	flagPollIntervalValue     = 10
	flagPriorityAliases       = "p"
	flagPriorityName          = "priority"
	flagPriorityUsage         = "a path to a file or directory to scan for executables"
	flagPriorityScriptAliases = "ps"
	flagPriorityScriptName    = "priority-script"
	flagPriorityScriptUsage   = "a path to a script that will run once when any priority starts"
	flagSwapAliases           = "s"
	flagSwapName              = "swap"
	flagSwapUsage             = "a process that will run when any priority executable is not running"
)

// NewApp returns a urfave/cli app that runs the loops to
// check for prioritized processes.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Action = run
	app.Authors = authors()
	app.Flags = flags()
	app.Name = appName
	app.Usage = appUsage
	app.UsageText = appUsageText

	return app
}

func authors() []*cli.Author {
	return []*cli.Author{
		{
			Name: authorName,
		},
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  flagDiableActionsName,
			Usage: flagDiableActionsUsage,
		},
		&cli.StringSliceFlag{
			Aliases: strings.Split(flagIgnoreAliases, ","),
			Name:    flagIgnoreName,
			Usage:   flagIgnoreUsage,
		},
		&cli.StringSliceFlag{
			Aliases:  strings.Split(flagPriorityAliases, ","),
			Name:     flagPriorityName,
			Usage:    flagPriorityUsage,
			Required: true,
		},
		&cli.StringFlag{
			Aliases: strings.Split(flagPriorityScriptAliases, ","),
			Name:    flagPriorityScriptName,
			Usage:   flagPriorityScriptUsage,
		},
		&cli.StringSliceFlag{
			Aliases:  strings.Split(flagSwapAliases, ","),
			Name:     flagSwapName,
			Usage:    flagSwapUsage,
			Required: true,
		},
		&cli.IntFlag{
			Aliases: strings.Split(flagLimitAliases, ","),
			Name:    flagLimitName,
			Usage:   flagLimitUsage,
			Value:   flagLimitValue,
		},
		&cli.IntFlag{
			Aliases: strings.Split(flagPollIntervalAliases, ","),
			Name:    flagPollIntervalName,
			Usage:   flagPollIntervalUsage,
			Value:   flagPollIntervalValue,
		},
	}
}

func run(c *cli.Context) error {
	loop := NewLoop()
	// Setup priority executables.
	pe := listExecutables(c.StringSlice(flagPriorityName), c.StringSlice(flagIgnoreName))
	if len(pe) == 0 {
		logWarn(fmt.Sprintf("%s found no priority executables - swap processes will run indefinitely", aurora.Cyan("setup")))
	} else {
		execs := strconv.Itoa(len(pe))
		logInfo(fmt.Sprintf("%s found %s priority executables", aurora.Cyan("setup"), aurora.Bold(execs)))
	}

	loop.WithPriorities(pe)

	// Setup swap scripts.
	// -s is a required flag, so there's no need to check if no swap processes
	// were passed in.
	sp := c.StringSlice(flagSwapName)
	s := swaps(pe, sp)
	loop.WithSwaps(s)

	swapCount := strconv.Itoa(len(sp))
	logInfo(fmt.Sprintf("%s registered %s swap processes", aurora.Cyan("setup"), aurora.Bold(swapCount)))
	// Setup priority script.
	// This will run once any priority starts. We wait for completion of the script.
	ps := c.String(flagPriorityScriptName)
	if ps != "" {
		loop.WithPriorityScript(ps)
		logInfo(fmt.Sprintf("%s registered priority script %s", aurora.Cyan("setup"), aurora.Bold(ps)))
	}
	// Set limit for loop to run.
	if limit := c.Int(flagLimitName); limit > 0 {
		loop.WithLimit(limit)
	}
	// Set the poll interval.
	pollInterval := c.Int(flagPollIntervalName)
	if pollInterval > 0 {
		loop.WithPollInterval(pollInterval)
	}
	// By default, enable all actions (keyboard inputs).
	if !c.Bool(flagDiableActionsName) {
		loop.WithActionsEnabled(true)
	}
	// This will run indefinitely unless limit is set to more than 0, or until the user exits.
	loop.Run()

	return nil
}

func listExecutables(paths, ignored []string) []*godirwalk.Dirent {
	// These are our "priority executables".
	pe := []*godirwalk.Dirent{}

	// Priority and swap process setup.
	for _, pd := range paths {
		e, err := ProcessList(pd, ignored)
		if err != nil {
			logError(fmt.Sprintf("%s error searching %s for executables: %s", aurora.Cyan("setup"), pd, err.Error()))

			continue
		}

		pe = append(pe, e...)
	}

	return pe
}

func swaps(pe []*godirwalk.Dirent, sp []string) []Swap {
	// Make sure there's no intersection here, that would be a nightmare.
	if err := intersect(pe, sp); err != nil {
		logFatal(err.Error())
	}

	swaps := []Swap{}
	for _, swap := range sp {
		swaps = append(swaps, NewSwap(swap))
	}

	return swaps
}

func intersect(files []*godirwalk.Dirent, swaps []string) error {
	filesMap := map[string]bool{}

	for _, file := range files {
		filesMap[file.Name()] = true
	}

	for _, swap := range swaps {
		file := filepath.Base(swap)
		if filesMap[file] {
			return fmt.Errorf("%s found in both priorities and swaps, this would be bad; exiting", file)
		}
	}

	return nil
}
