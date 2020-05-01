package servlet

import (
	"fmt"
	"time"
)

/// ---------------------------------------------------------------------------
/// TriggerCallback
/// ---------------------------------------------------------------------------

// TriggerCallback used as a trigger execution process.
type TriggerCallback func() error

/// ---------------------------------------------------------------------------
/// Trigger
/// ---------------------------------------------------------------------------

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

/// ---------------------------------------------------------------------------
/// PulseTrigger
/// ---------------------------------------------------------------------------

type PulseTrigger struct {
	Trigger
}

// NewPulseTrigger instantiate a new pulse trigger that will execute a
// callback method after a determined amount of time.
func NewPulseTrigger(delay time.Duration, callback TriggerCallback) (*PulseTrigger, error) {
	if callback == nil {
		return nil, fmt.Errorf("invalid nil 'callback' argument")
	}

	t := &PulseTrigger{
		Trigger: Trigger{
			timer:       delay,
			callback:    callback,
			isStopped:   false,
			channelStop: make(chan bool),
		},
	}

	go func() {
		for {
			select {
			case <-time.After(t.timer):
				if !t.isStopped {
					t.isStopped = true
					_ = t.callback()
				}
				return
			case <-t.channelStop:
				return
			}
		}
	}()

	return t, nil
}

/// ---------------------------------------------------------------------------
/// RecurringTrigger
/// ---------------------------------------------------------------------------

type RecurringTrigger struct {
	Trigger
}

// NewRecurringTrigger instantiate a new trigger that will execute a
// callback method recurrently with a defined periodicity.
func NewRecurringTrigger(period time.Duration, callback TriggerCallback) (*RecurringTrigger, error) {
	if callback == nil {
		return nil, fmt.Errorf("invalid nil 'callback' argument")
	}

	t := &RecurringTrigger{
		Trigger: Trigger{
			timer:       period,
			callback:    callback,
			isStopped:   false,
			channelStop: make(chan bool),
		},
	}

	go func() {
		for {
			select {
			case <-time.After(t.timer):
				if !t.isStopped {
					if err := t.callback(); err != nil {
						t.isStopped = true
						return
					}
				}
			case <-t.channelStop:
				t.isStopped = true
				return
			}
		}
	}()

	return t, nil
}
