package network

import "time"

// PeriodicTask represents a task that runs at regular intervals.
type PeriodicTask struct {
	ticker *time.Ticker
	stop   chan bool
}

// NewPeriodicTask creates and starts a new periodic task.
// The task runs immediately upon creation and then at the specified intervals.
// taskFunc is the function that will be executed periodically.
// Returns a pointer to the PeriodicTask.
func NewPeriodicTask(interval time.Duration, taskFunc func()) *PeriodicTask {
	pt := &PeriodicTask{
		ticker: time.NewTicker(interval),
		stop:   make(chan bool),
	}

	go func() {
		for {
			taskFunc()

			select {
			case <-pt.ticker.C:
				taskFunc()
			case <-pt.stop:
				pt.ticker.Stop()
				return
			}
		}
	}()

	return pt
}

// Stop terminates the periodic execution of the task.
// It stops the ticker and sends a signal to terminate the task's goroutine.
func (pt *PeriodicTask) Stop() {
	pt.stop <- true
}
