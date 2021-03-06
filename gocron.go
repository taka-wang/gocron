// Package gocron : A Golang Job Scheduling Package.
//
// An in-process scheduler for periodic jobs that uses the builder pattern
// for configuration. Schedule lets you run Golang functions periodically
// at pre-determined intervals using a simple, human-friendly syntax.
//
// Inspired by the Ruby module clockwork <https://github.com/tomykaira/clockwork>
// and
// Python package schedule <https://github.com/dbader/schedule>
//
// See also
// http://adam.heroku.com/past/2010/4/13/rethinking_cron/
// http://adam.heroku.com/past/2010/6/30/replace_cron_with_clockwork/
//
// maintained by Mark Salpeter mark@dealyze.com
//
// Copyright 2014 Jason Lyu. jasonlvhit@gmail.com .
// All rights reserved.
// Use of this source code is governed by a BSD-style .
// license that can be found in the LICENSE file.
package gocron

import (
	"time"
)

var defaultScheduler = NewScheduler()

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
//
// Please note that it is *intended behavior that `RunPending()`
// does not run missed jobs*. For example, if you've registered a job
// that should run every minute and you only call `RunPending()`
// in one hour increments then your job will only be run once every hour
func RunPending() {
	defaultScheduler.RunPending()
}

// RunAll runs all jobs of the regardless if they are scheduled to run or not.
// i.e., runs all jobs immediately
func RunAll() {
	defaultScheduler.RunAll()
}

// RunAllWithDelay runs all of the jobs with a delay between each of them
//
// This can help to distribute the system load generated by the jobs more evenly over
// time.
//
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
