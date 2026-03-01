package main

import (
	"flag"
	"log"
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

	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	app := Application{}
	defer app.Close()

	app.Init()
	app.Run()
}
