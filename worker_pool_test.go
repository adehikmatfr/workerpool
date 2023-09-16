package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	maxWorkers := 3
	wm := NewWorkerManager(maxWorkers)

	// Create a wait group to wait for all jobs to complete
	var wg sync.WaitGroup

	// Define a job that increments a counter
	counter := 0
	jobFunc := func() {
		counter++
		time.Sleep(100 * time.Millisecond)
		wg.Done() // Signal that the job is done
	}

	// Submit 5 jobs
	numJobs := 5
	wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		wm.SubmitJob(jobFunc)
	}

	// Stop the worker pool and wait for jobs to complete
	wm.StopAndWait()

	// Check if all jobs have been executed
	if counter != numJobs {
		t.Errorf("Expected %d jobs to be executed, but got %d", numJobs, counter)
	}
}

func TestWorkerPoolWithTimeout(t *testing.T) {
	maxWorkers := 2
	wm := NewWorkerManager(maxWorkers)

	// Define a job that sleeps for a while
	jobFunc := func() {
		time.Sleep(100 * time.Millisecond)
	}

	// Submit a job
	wm.SubmitJob(jobFunc)

	// Start a timer to measure how long it takes to stop the worker pool
	startTime := time.Now()
	wm.StopAndWait()
	elapsedTime := time.Since(startTime)

	// Check if the worker pool stops within a reasonable time
	if elapsedTime > 200*time.Millisecond {
		t.Errorf("Worker pool did not stop in a timely manner")
	}
}
