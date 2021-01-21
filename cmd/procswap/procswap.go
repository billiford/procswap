package main

import (
	"log"
	"os"

	"github.com/billiford/procswap/pkg/dir"
	"github.com/billiford/procswap/pkg/loop"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Aliases:  []string{"pd"},
				Name:     "priority-directories",
				Usage:    "a list of directories that will be scanned for executables",
				Required: true,
			},
			&cli.StringSliceFlag{
				Aliases:  []string{"bs"},
				Name:     "background-scripts",
				Usage:    "scripts that will run when any priority executable is not running",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			l := loop.New()

			// TODO add pretty logging.
			// Priority processes and background-scripts.
			var pe []os.FileInfo
			for _, pd := range c.StringSlice("priority-directories") {
				log.Printf("searching %s for executables\n", pd)

				e, err := dir.ListForExecutables(pd)
				if err != nil {
					return err
				}
				pe = append(pe, e...)
			}

			log.Println("total priority executables found:", len(pe))

			l.WithPriorityExecutables(pe)

			log.Printf("registering %d background scripts\n", len(c.StringSlice("background-scripts")))

			l.WithBackgroundScripts(c.StringSlice("background-scripts"))
			l.Run()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
