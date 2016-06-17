package main

import (
	"fmt"
	"time"

	"github.com/taka-wang/gocron"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func main() {
	// Do jobs with params
	gocron.Every(1).Second().Do(taskWithParams, 1, "hello")

	// Do jobs without params
	job1 := gocron.Every(1).Second().Do(task)
	gocron.Every(2).Seconds().Do(task)

	gocron.Remove(job1)
	// gocron.Clear()

	// function Start start all the pending jobs
	gocron.Start()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently

	for {
		time.Sleep(300 * time.Millisecond)
	}
}
