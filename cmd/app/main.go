package main

import (
	"flag"
	"os"
)

func main() {
	var (
		directory = flag.String("d", ".", "working directory")
	)
	flag.Parse()
	if *directory != "." {
		os.Chdir(*directory)
	}

	app := Application{}
	app.Init()
	app.Run()
}
