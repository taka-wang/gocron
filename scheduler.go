package gocron

import (
	"sort"
	"sync"
	"time"
)

// Scheduler keeps a slice of jobs that it executes at a regular interval
type Scheduler interface {

	// Clear removes all of the jobs that have been added to the scheduler
	Clear()

	// Emergency create a emergency job, and adds it to the `Scheduler`
	Emergency() *Job

	// Every creates a new job, and adds it to the `Scheduler`
	Every(interval uint64) *Job

	// EveryWithName creates a new job, and adds it to the `Scheduler` and job Map
	EveryWithName(interval uint64, name string) *Job

	// IsRunning returns true if the job  has started
	IsRunning() bool

	// Location sets the default location of every job created with `Every`.
	// The default location is `time.Local`
	Location(*time.Location)

	// NextRun returns the next next job to be run and the time in which
	// it will be run
	NextRun() (*Job, time.Time)

	// Remove removes an individual job from the scheduler. It returns true
	// if the job was found and removed from the `Scheduler`
	Remove(*Job) bool

	// UpdateIntervalWithName update an individual job's interval from the scheduler by name.
	// It returns true if the job was found and update interval
	UpdateIntervalWithName(name string, interval uint64) bool

	// RemoveWithName removes an individual job from the scheduler by name. It returns true
	// if the job was found and removed from the `Scheduler`
	RemoveWithName(string) bool

	// PauseWithName pause an individual job by name. It returns true if the job was found and set enabled
	PauseWithName(string) bool

	// PauseAll disable all jobs
	PauseAll()

	// ResumeWithName resume an individual job by name. It returns true if the job was found and set enabled
	ResumeWithName(string) bool

	// ResumeAll resume all jobs
	ResumeAll()

	// Depricated: RunAll runs all of the jobs regardless of wether or not
	// they are pending
	RunAll()

	// RunAllWithDelay runs all of the jobs regardless of wether or not
	// they are pending with a delay
	RunAllWithDelay(time.Duration)

	// Depricated: RunPending runs all of the pending jobs
	RunPending()

	// Start starts the scheduler
	Start()

	// Stop stops the scheduler from executing jobs
	Stop()
}

// NewScheduler create a new scheduler.
// Note: the current implementation is not concurrency safe.
func NewScheduler() Scheduler {
	return &scheduler{
		jobMap:    make(map[string]*Job),
		isStopped: make(chan bool),
		location:  time.Local,
	}
}

// Scheduler contains jobs and a loop to run the jobs
type scheduler struct {
	jobMap    map[string]*Job
	ejobs     []*Job // Emergency jobs
	jobs      []*Job
	isRunning bool
	isStopped chan bool
	location  *time.Location
	mutex     sync.Mutex
}

// Len returns the number of jobs that have been scheduled.
// It is part of the `sort.Interface` interface
func (s *scheduler) Len() int {
	return len(s.jobs)
}

// Swap swaps the order of two jobs in the backing slice.
// It is part of the `sort.Interface` interface
func (s *scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

// Less swaps the order of two jobs in the backing slice.
// It is part of the `sort.Interface` interface
func (s *scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.After(s.jobs[i].nextRun)
}

// NextRun returns the job and time when the next job should run
func (s *scheduler) NextRun() (*Job, time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.jobs) == 0 {
		return nil, time.Time{}
	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}

// Every schedules a new job
func (s *scheduler) Every(interval uint64) *Job {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	job := newJob(interval).Location(s.location)
	s.jobs = append(s.jobs, job)

	// inserion sort by 'interval'
	for i := 1; i < len(s.jobs); i++ {
		tmp := s.jobs[i]
		j := i - 1
		for j >= 0 && s.jobs[j].interval > tmp.interval {
			s.jobs[j+1] = s.jobs[j]
			j = j - 1
		}
		s.jobs[j+1] = tmp
	}

	return job
}

// Add job name and job object to jobMap
func (s *scheduler) EveryWithName(interval uint64, name string) *Job {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// if job exist, remove it;
	if oldJob, ok := s.jobMap[name]; ok {
		// we don't call s.Remove since it cause deadlock
		for i, job := range s.jobs {
			if oldJob == job {
				copy(s.jobs[i:], s.jobs[i+1:])
				s.jobs[len(s.jobs)-1] = nil
				s.jobs = s.jobs[:len(s.jobs)-1]
			}
		}
	}

	// create/update job to job list and job map
	job := newJob(interval).Location(s.location)
	s.jobMap[name] = job
	s.jobs = append(s.jobs, job)

	// inserion sort by 'interval'
	for i := 1; i < len(s.jobs); i++ {
		tmp := s.jobs[i]
		j := i - 1
		for j >= 0 && s.jobs[j].interval > tmp.interval {
			s.jobs[j+1] = s.jobs[j]
			j = j - 1
		}
		s.jobs[j+1] = tmp
	}

	return job
}

// Emergency schedules a new emergency job
func (s *scheduler) Emergency() *Job {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// cheat the interval
	job := newJob(1).Location(s.location)
	s.ejobs = append(s.ejobs, job)

	return job
}

// runPending runs all of the jobs pending at this time
func (s *scheduler) runPending(now time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// run emergency jobs
	for _, job := range s.ejobs {
		job.run()
	}
	// clear ejobs queue
	s.ejobs = []*Job{}

	sort.Sort(s)
	// run jobs
	for _, job := range s.jobs {
		if !job.isInit() {
			// set lastRun and nextRun
			job.init(now)
		}
		if job.shouldRun(now) {
			job.run()
		} else {
			// intend to loop through
			continue
		}
	}
}

// Depricated: RunPending runs all of the jobs that are scheduled to run
func (s *scheduler) RunPending() {
	s.runPending(time.Now())
}

// Depricated: RunAll rungs all jobs regardless if they are scheduled to run or not
func (s *scheduler) RunAll() {
	RunAllWithDelay(0)
}

// Depricated: RunAllWithDelay all jobs with delay seconds
func (s *scheduler) RunAllWithDelay(d time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	sort.Sort(s)
	for _, job := range s.jobs {
		if !job.isInit() {
			// set lastRun and nextRun
			job.init(now)
		}
		// force to run
		job.run()
		time.Sleep(d)
	}

}

// Location sets the default location for every job created
// with `Scheduler.Every(...)`. By default the location is `time.Local`
func (s *scheduler) Location(location *time.Location) {
	s.location = location
}

// Removes a job from the queue
func (s *scheduler) Remove(j *Job) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, job := range s.jobs {
		if j == job {
			// fix potential memory leak problem arrcording to:
			// https://github.com/golang/go/wiki/SliceTricks
			copy(s.jobs[i:], s.jobs[i+1:])
			s.jobs[len(s.jobs)-1] = nil
			s.jobs = s.jobs[:len(s.jobs)-1]
			return true
		}
	}

	return false
}

// RemoveWithName removes an individual job from the scheduler by name
func (s *scheduler) RemoveWithName(name string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if dJob, ok := s.jobMap[name]; ok {
		// we don't call s.Remove since it cause deadlock
		for i, job := range s.jobs {
			if dJob == job {
				copy(s.jobs[i:], s.jobs[i+1:])
				s.jobs[len(s.jobs)-1] = nil
				s.jobs = s.jobs[:len(s.jobs)-1]
				delete(s.jobMap, name) // remove jobMap item
				return true
			}
		}
	}
	return false
}

// UpdateIntervalWithName  update interval by name
func (s *scheduler) UpdateIntervalWithName(name string, interval uint64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, ok := s.jobMap[name]; ok {
		job.updateInterval(interval)
		return true
	}
	return false
}

// PauseWithName disable job by name
func (s *scheduler) PauseWithName(name string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, ok := s.jobMap[name]; ok {
		job.pause()
		return true
	}
	return false
}

// PauseAll disable all jobs
func (s *scheduler) PauseAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range s.jobMap {
		v.pause()
	}
}

// ResumeWithName enable job by name
func (s *scheduler) ResumeWithName(name string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if job, ok := s.jobMap[name]; ok {
		job.resume()
		return true
	}
	return false
}

// ResumeAll enable all jobs
func (s *scheduler) ResumeAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, v := range s.jobMap {
		v.resume()
	}
}

// Clear deletes all scheduled jobs
func (s *scheduler) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.jobs = []*Job{}
	s.jobMap = make(map[string]*Job) // new job map
}

// Start all the pending jobs
// Add seconds ticker
func (s *scheduler) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// only start the scheduler if it hasn't been started yet
	if s.isRunning {
		return
	}

	// start the scheduler
	isStarted := make(chan bool)
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for {
			select {
			case now := <-ticker.C:
				if !s.isRunning {
					// initialize all of the jobs with the first ticker time
					// so that they are all in sync with the run loop
					for _, job := range s.jobs {
						job.init(now)
					}
					s.isRunning = true
					isStarted <- true
				}
				s.runPending(now)
			case <-s.isStopped:
				s.isRunning = false
				// send a confirmation message back to the `Stop()` method
				s.isStopped <- true
				return
			}
		}
	}()

	// wait until he ticker has been started and all of the jobs
	// have been initialized
	<-isStarted
}

// IsRunning returns true if the scheduler is startes
func (s *scheduler) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.isRunning
}

// Stop stops the scheduler
func (s *scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// only send the stop signal if the scheduler has been started
	if s.isRunning {
		s.isStopped <- true
		// wait for the ticker to send a confirmation message back through
		// the stop channel just before it shuts down the ticker loop
		<-s.isStopped
	}
}
