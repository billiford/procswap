package procswap

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli/v2"
)

const (
	flagPriorityAliases = "p"
	flagPriorityName    = "priority"
	flagPriorityUsage   = "a path to a file or directory to scan for executables"
	flagSwapAliases     = "s"
	flagSwapName        = "swap"
	flagSwapUsage       = "a process that will run when any priority executable is not running"
)

// NewApp returns a urfave/cli app that runs the loops to
// check for prioritized processes.
func NewApp() *cli.App {
	return &cli.App{
		Name:  "procswap",
		Usage: "prioritize processes",
		// Version:     "",
		Description: "run processes when any prioritized process is not running",
		Commands:    []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Aliases:  strings.Split(flagPriorityAliases, ","),
				Name:     flagPriorityName,
				Usage:    flagPriorityUsage,
				Required: true,
			},
			&cli.StringSliceFlag{
				Aliases:  strings.Split(flagSwapAliases, ","),
				Name:     flagSwapName,
				Usage:    flagSwapUsage,
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			run(c)

			return nil
		},
		Authors: []*cli.Author{
			&cli.Author{
				Name: "billiford",
			},
		},
	}
}

func run(c *cli.Context) {
	loop := NewLoop()
	// These are our "priority executables".
	pe := []os.FileInfo{}

	// Priority and swap process setup.
	paths := c.StringSlice(flagPriorityName)
	for _, pd := range paths {
		e, err := ProcessList(pd)
		if err != nil {
			LogError(fmt.Sprintf("error searching %s for executables: %s", pd, err.Error()))

			continue
		}

		pe = append(pe, e...)
	}

	if len(pe) == 0 {
		LogWarn("found no priority executables - swap processes will run indefinitely")
	} else {
		execs := strconv.Itoa(len(pe))
		LogInfo(fmt.Sprintf("found %s priority executables", aurora.Bold(execs)))
	}

	loop.WithPriorities(pe)

	// -sp is a required flag, so there's no need to check if no swap processes
	// were passed in.
	sp := c.StringSlice(flagSwapName)
	loop.WithSwaps(sp)

	swaps := strconv.Itoa(len(sp))
	LogInfo(fmt.Sprintf("registered %s swap processes", aurora.Bold(swaps)))

	// This will run indefinitely, until the user exits.
	loop.Run()
}
