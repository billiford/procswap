package main

import (
	"log"
	"os"

	procswap "github.com/billiford/procswap/pkg"
	"github.com/mattn/go-colorable"
)

func main() {
	// Disable the log package from printing the date and time - we will handle that.
	log.SetFlags(0)
	// Need to set this since displaying color on a Windows console is tough.
	log.SetOutput(colorable.NewColorableStdout())

	// Create a new app. This is a urfave/cli app making it easier to setup.
	app := procswap.NewApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
