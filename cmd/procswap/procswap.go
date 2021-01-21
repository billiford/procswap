package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	procswap "github.com/billiford/procswap/pkg"
	"github.com/billiford/procswap/pkg/dir"
	"github.com/billiford/procswap/pkg/loop"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli/v2"

	"github.com/logrusorgru/aurora"
)

const (
	flagPriorityAliases = "p"
	flagPriorityName    = "priority"
	flagPriorityUsage   = "a directory to scan for executables"
	flagSwapAliases     = "s"
	flagSwapName        = "swap"
	flagSwapUsage       = "a process that will run when any priority executable is not running"
)

func main() {
	// Disable the log package from printing the date and time - we will handle that.
	log.SetFlags(0)
	log.SetOutput(colorable.NewColorableStdout())

	app := &cli.App{
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
			l := loop.New()

			// Priority and swap process setup.
			var pe []os.FileInfo
			for _, pd := range c.StringSlice(flagPriorityName) {
				procswap.LogInfo(fmt.Sprintf("searching %s for executables", pd))

				e, err := dir.ListForExecutables(pd)
				if err != nil {
					procswap.LogError(fmt.Sprintf("error searching %s for executables: %s", pd, err.Error()))

					continue
				}
				pe = append(pe, e...)
			}

			if len(pe) == 0 {
				procswap.LogWarn("found no priority executables - swap processes will run indefinitely")
			} else {
				execs := strconv.Itoa(len(pe))
				procswap.LogInfo(fmt.Sprintf("found %s priority executables", aurora.Bold(execs)))
			}

			l.WithPriorities(pe)

			sp := c.StringSlice(flagSwapName)
			if len(sp) == 0 {
				procswap.LogFatal("no swap processes passed in")
			}

			l.WithSwaps(sp)

			swaps := strconv.Itoa(len(sp))
			procswap.LogInfo(fmt.Sprintf("registered %s swap processes", aurora.Bold(swaps)))

			l.Run()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
