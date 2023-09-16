package workerpool

import (
	"sync"
)

type Worker interface {
	start()
	stop()
}

type worker struct {
	jobQueue   chan Job
	stopSignal chan struct{}
	workerNum  int // Worker number
	wg         *sync.WaitGroup
}

func newWorker(job chan Job, signal chan struct{}, wg *sync.WaitGroup, num int) Worker {
	return &worker{
		jobQueue:   job,
		stopSignal: signal,
		workerNum:  num,
		wg:         wg,
	}
}

func (w *worker) start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case job, ok := <-w.jobQueue:
				if !ok {
					return // Worker was signaled to stop
				}
				job() // Execute the job
				// log here if you want to know the job that has been executed and the worker number
				// fmt.Printf("Job Finished By Worker %d \n", w.workerNum)
			case <-w.stopSignal:
				return // Worker was signaled to stop
			}
		}
	}()
}

func (w *worker) stop() {
	select {
	case <-w.stopSignal:
		// worker has been executed
		// fmt.Printf("Worker %d Close signal \n", w.workerNum)
	default:
		close(w.stopSignal)
	}
}
