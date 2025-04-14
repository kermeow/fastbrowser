package main

import (
	"fastgh3/fastbrowser/browser"
	"fastgh3/fastbrowser/config"
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	conf, _ := config.Load()

	ui := browser.New(conf)
	go func() {
		err := ui.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
