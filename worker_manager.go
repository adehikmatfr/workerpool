package workerpool

import "sync"

type Job func()

// WorkerManager represents a pool of worker goroutines that can execute jobs.
type workerManager struct {
	workers    []Worker
	jobQueue   chan Job
	stopSignal chan struct{}
	closeOnce  sync.Once
	wg         *sync.WaitGroup
}

type WorkerManager interface {
	SubmitJob(job Job)
	StopAndWait()
}

func NewWorkerManager(maxWorker int) WorkerManager {
	if maxWorker < 1 {
		maxWorker = 1
	}

	// create worker manager
	wm := &workerManager{
		jobQueue:   make(chan Job),
		stopSignal: make(chan struct{}),
		wg:         &sync.WaitGroup{},
	}

	for i := 0; i < maxWorker; i++ {
		worker := newWorker(wm.jobQueue, wm.stopSignal, wm.wg, i)
		worker.start() // start worker
		wm.workers = append(wm.workers, worker)
	}

	return wm
}

// SubmitJob adds a job to the worker pool's job queue.
func (wm *workerManager) SubmitJob(job Job) {
	wm.jobQueue <- job
}

// StopAndWait stops all workers and waits for them to finish.
func (wm *workerManager) StopAndWait() {
	wm.closeOnce.Do(func() {
		close(wm.jobQueue) // Close the job queue once using sync.Once
		close(wm.stopSignal)
	})

	for _, worker := range wm.workers {
		worker.stop()
	}

	wm.wg.Wait()
}
