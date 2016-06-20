## goCron: A Golang Job Scheduling Package.
[![Build Status](https://travis-ci.org/taka-wang/gocron.svg?branch=master)](https://travis-ci.org/taka-wang/gocron)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/taka-wang/gocron)

goCron is a Golang job scheduling package which lets you run Go functions periodically at pre-determined interval using a simple, human-friendly syntax.

goCron is a Golang implementation of Ruby module [clockwork](<https://github.com/tomykaira/clockwork>) and Python job scheduling package [schedule](<https://github.com/dbader/schedule>), and personally, this package is my first Golang program, just for fun and practice.

See also this two great articles:
* [Rethinking Cron](http://adam.heroku.com/past/2010/4/13/rethinking_cron/)
* [Replace Cron with Clockwork](http://adam.heroku.com/past/2010/6/30/replace_cron_with_clockwork/)

Back to this package, you could just use this simple API as below, to run a cron scheduler.

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

## APIs
```go
// Every schedules a new job in the default scheduler
func Every(interval uint64) *Job {
	return defaultScheduler.Every(interval)
}

// EveryWithName schedules a new job in the default scheduler
func EveryWithName(interval uint64, name string) *Job {
	return defaultScheduler.EveryWithName(interval, name)
}

// Emergency schedules a new emergency job in the default scheduler
func Emergency() *Job {
	return defaultScheduler.Emergency()
}

// RunPending Runs all of the jobs that are scheduled to run
func RunPending() {
	defaultScheduler.RunPending()
}

// RunAll runs all jobs of the regardless if they are scheduled to run or not.
// i.e., runs all jobs immediately
func RunAll() {
	defaultScheduler.RunAll()
}

// RunAllWithDelay runs all of the jobs with a delay between each of them
func RunAllWithDelay(d time.Duration) {
	defaultScheduler.RunAllWithDelay(d)
}

// Start starts the scheduler
func Start() {
	defaultScheduler.Start()
}

// IsRunning returns true if the default scheduler has started
func IsRunning() bool {
	return defaultScheduler.IsRunning()
}

// Stop stops the default scheduler
func Stop() {
	defaultScheduler.Stop()
}

// Clear removes all of the jobs from the default scheduler
func Clear() {
	defaultScheduler.Clear()
}

// Remove removes the job from the default scheduler
func Remove(j *Job) {
	defaultScheduler.Remove(j)
}

// RemoveWithName removes an individual job from the default scheduler
func RemoveWithName(name string) {
	defaultScheduler.RemoveWithName(name)
}

// PauseWithName pause an individual job by name from the default scheduler
func PauseWithName(name string) {
	defaultScheduler.PauseWithName(name)
}

// ResumeWithName resume an individual job by name from the default scheduler
func ResumeWithName(name string) {
	defaultScheduler.ResumeWithName(name)
}

// NextRun gets the next running time
func NextRun() (job *Job, time time.Time) {
	return defaultScheduler.NextRun()
}

```
