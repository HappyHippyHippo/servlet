package servlet

import (
	"fmt"
	"time"
)

/// ---------------------------------------------------------------------------
/// Trigger
/// ---------------------------------------------------------------------------

// Trigger defines the base trigger instance functionality.
type Trigger struct {
	timer       time.Duration
	callback    TriggerCallback
	isStopped   bool
	channelStop chan bool
}

// Close will retrieve the time period associated to the trigger.
func (t *Trigger) Close() {
	if t == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	t.Stop()
}

// Timer will retrieve the time period associated to the trigger.
func (t Trigger) Timer() time.Duration {
	return t.timer
}

// IsStopped check if the trigger is stopped.
func (t Trigger) IsStopped() bool {
	return t.isStopped
}

// Stop signals the trigger to stop execution.
func (t *Trigger) Stop() {
	if t == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if !t.isStopped {
		t.isStopped = true
		t.channelStop <- true
	}
}
