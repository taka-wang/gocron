// Tests for gocron
package gocron

import (
	"fmt"
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
)

func task() {
	fmt.Println("I am runnning task.")
}

func task2() {
	fmt.Println("I am runnning task.")
}

func task3() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	t := time.Now()
	fmt.Println(a, b, t.Format("2006-01-02 15:04:05.000"))
}

/*
func TestJob(t *testing.T) {

	// note: we're defining today as the first of the month so we can test an important edge case in the lastRun
	// calculation when a jobs lastRun time occured the previous month from when the job initialized
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), 1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	aMinuteAgo := today.Add(-time.Minute)
	aMinuteFromNow := today.Add(time.Minute)
	aMinuteAgoAtTime := fmt.Sprintf("%02d:%02d", aMinuteAgo.Hour(), aMinuteAgo.Minute())
	aMinuteFromNowAtTime := fmt.Sprintf("%02d:%02d", aMinuteFromNow.Hour(), aMinuteFromNow.Minute())
	s := sugar.New(t)

	s.Title("Day")

	s.Assert("`Job.Every(...).Day().At(...)`", func(log sugar.Log) bool {
		// try this with 20 random day intervals
		for i := 20; i > 0; i-- {

			// get a random interval of days [1, 5]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%6

			// create and init the job
			job := newJob(interval).Day().At(aMinuteFromNowAtTime)
			job.init(today)

			// jobs last run should be `interval` days from now
			aMinuteFromNowIntervalDaysFromnow := aMinuteFromNow
			if !job.lastRun.Equal(aMinuteFromNow) {
				log("the lastRun did not occur a minute from now days ago")
				log(job.lastRun, aMinuteFromNowIntervalDaysFromnow)
				return false
			}
			//
			// // jobs last run should be `interval` days from now
			// aMinuteAgoIntervalDaysAgo := aMinuteAgo.Add(-1 * Day * time.Duration(interval))
			// if !job.lastRun.Equal(aMinuteAgoIntervalDaysAgo) {
			// 	log("the lastRun did not occur %d days ago", interval)
			// 	log(job.lastRun, aMinuteFromNowIntervalDaysAgo)
			// 	return false
			// }

			// jobs next run is should be today
			if !job.nextRun.Equal(aMinuteFromNow) {
				log("the nextRun will not happen a minute from now")
				log(job.nextRun, aMinuteFromNow)
				return false
			}

			// after run, the nextRun is interval days from the previous nextRun
			job.run()
			aMinutFromNowIntervalDaysAfterNextRun := aMinuteFromNow.Add(Day * time.Duration(interval))
			if !job.nextRun.Equal(aMinutFromNowIntervalDaysAfterNextRun) {
				log("the next nextRun will not happen in %d days", interval)
				log(job.nextRun, aMinutFromNowIntervalDaysAfterNextRun)
				return false
			}

		}

		return true
	})

	s.Assert("`Job.Every(...).Day.At(...)` set to the past", func(log sugar.Log) bool {
		// try this with 20 random day intervals
		for i := 20; i > 0; i-- {

			// get a random interval of days [1, 5]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%5

			// create and init the job
			job := newJob(interval).Day().At(aMinuteAgoAtTime)
			job.init(today)

			// jobs last run interval days from tomorrow
			aMinuteAgoIntervalDaysFromTomorrow := aMinuteAgo.Add(Day).Add(-1 * Day * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteAgoIntervalDaysFromTomorrow) {
				log("the lastRun did not %d days from tomorrow", interval)
				log(job.lastRun, aMinuteAgoIntervalDaysFromTomorrow)
				return false
			}

			// jobs next run is tomorrow
			aMinuteAgoTomorrow := aMinuteAgo.Add(Day)
			if !job.nextRun.Equal(aMinuteAgoTomorrow) {
				log("the nextRun will not occur tomorrow")
				log(job.nextRun, aMinuteAgoTomorrow)
				return false
			}

			// after run, the nextRun is interval days from the previous nextRun
			job.run()
			aMinutAgoIntervalDaysAfterNextRun := aMinuteAgoTomorrow.Add(Day * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalDaysAfterNextRun) {
				log("the next nextRun will not happen in %d days", interval)
				log(job.nextRun, aMinutAgoIntervalDaysAfterNextRun)
				return false
			}
		}

		return true
	})

	s.Title("Week")

	s.Assert("`Job.Every(...).Weekday(...).At(...)` set to the past", func(log sugar.Log) bool {
		// try this with 20 random weekdays and week intervals
		for i := 20; i > 0; i-- {

			// get a random interval of weeks [1, 52]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%52

			// get a random day of the week that is today or before today
			rand.Seed(time.Now().UnixNano())
			weekday := time.Weekday(rand.Int() % int(today.Weekday()+1))
			durationAfterWeekday := time.Duration(weekday-today.Weekday()) * 24 * time.Hour

			// create and init the job
			job := newJob(interval).Weekday(weekday).At(aMinuteAgoAtTime)
			job.init(today)

			// jobs lastRun was interval weeks ago from next week
			aMinuteAgoIntervalWeeksFromNextWeek := aMinuteAgo.Add(durationAfterWeekday).Add(Week).Add(-1 * Week * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteAgoIntervalWeeksFromNextWeek) {
				log("the lastRun did not occur %d weeks ago", interval+1)
				log(weekday, aMinuteAgoIntervalWeeksFromNextWeek.Weekday(), job.nextRun.Weekday(), job.lastRun, aMinuteAgoIntervalWeeksFromNextWeek)
				return false
			}

			// jobs next run is next week
			aMinuteAgoNextWeek := aMinuteAgo.Add(durationAfterWeekday).Add(Week)
			if !job.nextRun.Equal(aMinuteAgoNextWeek) {
				log("the nextRun will not occur next week")
				log(weekday, aMinuteAgoNextWeek.Weekday(), job.nextRun.Weekday(), job.nextRun, aMinuteAgoNextWeek)
				return false
			}

			// after run, the nextRun is interval weeks from the previous nextRun
			job.run()
			aMinutAgoIntervalWeeksAfterNextRun := aMinuteAgoNextWeek.Add(Week * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalWeeksAfterNextRun) {
				log("the next nextRun will not happen in %d weeks", interval)
				log(job.nextRun, aMinutAgoIntervalWeeksAfterNextRun)
				return false
			}

		}

		return true
	})

	s.Assert("`Job.Every(...).Weekday(...).At(...)` set to the future", func(log sugar.Log) bool {
		// try this with 20 random weekdays and week intervals
		for i := 20; i > 0; i-- {

			// get a random interval of weeks [1, 52]
			rand.Seed(time.Now().UnixNano())
			interval := 1 + uint64(rand.Int())%52

			// get a random day of the week that is today or after today
			rand.Seed(time.Now().UnixNano())
			weekday := time.Weekday(int(today.Weekday()) + rand.Int()%(7-int(today.Weekday())))
			durationUntilWeekday := time.Duration(weekday-today.Weekday()) * 24 * time.Hour

			// create and init the job
			job := newJob(interval).Weekday(weekday).At(aMinuteFromNowAtTime)
			job.init(today)

			// jobs last run was interval weeks ago
			aMinuteFromNowIntervalWeeksAgo := aMinuteFromNow.Add(durationUntilWeekday).Add(-1 * Week * time.Duration(interval))
			if !job.lastRun.Equal(aMinuteFromNowIntervalWeeksAgo) {
				log("the lastRun did not occur %d weeks ago", interval)
				log(weekday, aMinuteFromNowIntervalWeeksAgo.Weekday(), job.nextRun.Weekday(), job.lastRun, aMinuteFromNowIntervalWeeksAgo)
				return false
			}

			// jobs next run is this week
			thisWeekdayAMinuteFromNow := aMinuteFromNow.Add(durationUntilWeekday)
			if !job.nextRun.Equal(thisWeekdayAMinuteFromNow) {
				log("the nextRun will not occur this week")
				log(weekday, thisWeekdayAMinuteFromNow.Weekday(), job.nextRun.Weekday(), job.nextRun, thisWeekdayAMinuteFromNow)
				return false
			}

			// after run, the nextRun is interval weeks from the previous nextRun
			job.run()
			aMinutAgoIntervalWeeksAfterNextRun := thisWeekdayAMinuteFromNow.Add(Week * time.Duration(interval))
			if !job.nextRun.Equal(aMinutAgoIntervalWeeksAfterNextRun) {
				log("the next nextRun will not happen in %d weeks", interval)
				log(job.nextRun, aMinutAgoIntervalWeeksAfterNextRun)
				return false
			}
		}

		return true
	})

	s.Title("Time")

	s.Assert("`Job.Hour()` causes lastRun to be now and nextRun to be `interval` hour(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})

	s.Assert("`Job.Minute()` causes lastRun to be now and nextRun to be `interval` minute(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})

	s.Assert("`Job.Second()` causes lastRun to be now and nextRun to be `interval` second(s) from now", func(log sugar.Log) bool {
		// TODO: implement test
		return false
	})
}

*/

func TestScheduler(t *testing.T) {

	s := sugar.New(t)

	/*
		s.Assert("`runPending(...)` runs all pending jobs", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})

		s.Assert("`Start()`, `IsRunning()` and `Stop()` perform correctly in asynchrnous environments", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})

		s.Assert("`Start()` triggers runPending(...) every second", func(log sugar.Log) bool {
			// TODO: implement test
			return false
		})
	*/

	s.Title("Job order test")

	s.Assert("`Every()` should append job with order", func(log sugar.Log) bool {

		s := scheduler{
			isStopped: make(chan bool),
			location:  time.Local,
		}

		s.Every(3).Seconds().Do(taskWithParams, 1, "3s")
		s.Every(2).Seconds().Do(taskWithParams, 2, "2s")
		s.Every(5).Seconds().Do(taskWithParams, 3, "5s")
		s.Every(1).Seconds().Do(taskWithParams, 4, "1s")
		s.Every(1).Seconds().Do(taskWithParams, 5, "1s")
		s.Every(500).Seconds().Do(taskWithParams, 6, "500s")
		s.Every(10).Seconds().Do(taskWithParams, 7, "10s")

		/*
			s.Start()

			time.Sleep(5 * time.Second)

			s.Every(1).Seconds().Do(taskWithParams, 8, "1s")
			time.Sleep(10 * time.Second)
		*/

		// debug
		for _, job := range s.jobs {
			log("interval: %d", job.interval)
		}

		if s.jobs[5].interval == 10 {
			return true
		}
		return false
	})

	s.Assert("`Remove()` should delete desired job", func(log sugar.Log) bool {
		s := scheduler{
			isStopped: make(chan bool),
			location:  time.Local,
		}

		// add three jobs
		s.Every(3).Seconds().Do(task)
		item := s.Every(2).Seconds().Do(task2)
		s.Every(1).Seconds().Do(task3)

		// debug
		for _, job := range s.jobs {
			log("interval: %d", job.interval)
		}

		// remove one job
		s.Remove(item)

		// debug
		for _, job := range s.jobs {
			log("@interval: %d", job.interval)
		}

		if s.Len() == 2 {
			return true
		}
		return false
	})

}
