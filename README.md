## goCron: A Golang Job Scheduling Package.

[![Build Status](http://dds.cmwang.net/api/badges/taka-wang/gocron/status.svg)](http://dds.cmwang.net/taka-wang/gocron)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/taka-wang/gocron)

goCron is a Golang job scheduling package which lets you run Go functions periodically at pre-determined interval using a simple, human-friendly syntax.



``` go
package main

import (
	"fmt"
	"github.com/taka-wang/gocron"
	"time"
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
	mjob := gocron.Every(1).Second().Do(task)
	gocron.Every(2).Seconds().Do(task)
	gocron.Every(1).Minute().Do(task)
	gocron.Every(2).Minutes().Do(task)
	gocron.Every(1).Hour().Do(task)
	gocron.Every(2).Hours().Do(task)
	gocron.Every(1).Day().Do(task)
	gocron.Every(2).Days().Do(task)

	// Do jobs on specific weekday
	gocron.Every(1).Monday().Do(task)
	gocron.Every(1).Thursday().Do(task)

	// function At() take a string like 'hour:min'
	gocron.Every(1).Day().At("10:30").Do(task)
	gocron.Every(1).Monday().At("18:30").Do(task)

	// remove, clear and next_run
	_, time := gocron.NextRun()
	fmt.Println(time)

	gocron.Remove(mjob)
	gocron.Clear()

	// function Start start all the pending jobs
	gocron.Start()
	
	// trigger emergency job
	gocron.Emergency().Do(taskWithParams, 9, "emergency")

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := gocron.NewScheduler()
	s.Every(3).Seconds().Do(task)
	s.Start()
	for {
		time.Sleep(300 * time.Millisecond)
	}

}
```
