package main

import (
	"os"

	"github.com/jasonlvhit/gocron"
)

var counter = 0

func main() {
	gocron.Every(15).Minute().Do(taskWithParams)
	startJob()
}

func startJob() {
	<-gocron.Start()
}

func stopJob() {
	gocron.Remove(taskWithParams)
	gocron.Clear()
	os.Exit(1)
}

func taskWithParams() {
	counter++
	if counter > 30 {
		stopJob()
	} else {
		// TODO: send email to pondthaitay@hotmail.com
	}
}
