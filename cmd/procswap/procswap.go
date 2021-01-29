package main

import (
	"fmt"
	"log"
	"os"

	procswap "github.com/billiford/procswap/pkg"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli/v2"
)

const (
	banner = `
      ___         ___           ___           ___           ___           ___           ___           ___
     /  /\       /  /\         /  /\         /  /\         /  /\         /__/\         /  /\         /  /\
    /  /::\     /  /::\       /  /::\       /  /:/        /  /:/_       _\_ \:\       /  /::\       /  /::\
   /  /:/\:\   /  /:/\:\     /  /:/\:\     /  /:/        /  /:/ /\     /__/\ \:\     /  /:/\:\     /  /:/\:\
  /  /:/~/:/  /  /:/~/:/    /  /:/  \:\   /  /:/  ___   /  /:/ /::\   _\_ \:\ \:\   /  /:/~/::\   /  /:/~/:/
 /__/:/ /:/  /__/:/ /:/___ /__/:/ \__\:\ /__/:/  /  /\ /__/:/ /:/\:\ /__/\ \:\ \:\ /__/:/ /:/\:\ /__/:/ /:/
 \  \:\/:/   \  \:\/:::::/ \  \:\ /  /:/ \  \:\ /  /:/ \  \:\/:/~/:/ \  \:\ \:\/:/ \  \:\/:/__\/ \  \:\/:/
  \  \::/     \  \::/~~~~   \  \:\  /:/   \  \:\  /:/   \  \::/ /:/   \  \:\ \::/   \  \::/       \  \::/
   \  \:\      \  \:\        \  \:\/:/     \  \:\/:/     \__\/ /:/     \  \:\/:/     \  \:\        \  \:\
    \  \:\      \  \:\        \  \::/       \  \::/        /__/:/       \  \::/       \  \:\        \  \:\
     \__\/       \__\/         \__\/         \__\/         \__\/         \__\/         \__\/         \__\/
`
)

var (
	version  string
	revision string
)

func main() {
	// Disable the log package from printing the date and time - we will handle that.
	log.SetFlags(0)
	// Need to set this since displaying color on a Windows console is tough.
	log.SetOutput(colorable.NewColorableStdout())

	cli.VersionPrinter = versionPrinter

	// Create a new app. This is a urfave/cli app making it easier to setup.
	app := procswap.NewApp()
	app.Version = version

	printBanner()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func versionPrinter(c *cli.Context) {
	log.Printf("version=%s revision=%s", c.App.Version, revision)
}

func printBanner() {
	fmt.Println(banner)
}
