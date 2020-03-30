package trigger

import (
	"time"
)

// Callback used as a trigger execution process.
type Callback func() error

// Trigger interface defines the interaction with a trigger.
type Trigger interface {
	Close() error
	Timer() time.Duration
	IsStopped() bool
	Stop() error
}

type trigger struct {
	timer       time.Duration
	callback    Callback
	isStopped   bool
	channelStop chan bool
}

// Timer will retrieve the time period associated to the trigger.
func (t *trigger) Timer() time.Duration {
	return t.timer
}

// IsStopped signal if the trigger is stopped.
func (t *trigger) IsStopped() bool {
	return t.isStopped
}

// Stop signal the trigger to stop execution.
func (t *trigger) Stop() error {
	if !t.isStopped {
		t.isStopped = true
		t.channelStop <- true
	}
	return nil
}
