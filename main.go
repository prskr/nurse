package main

import (
	"fmt"
	"os"

	"code.icb4dc0.de/prskr/nurse/cmd"
)

func main() {
	app, err := cmd.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
