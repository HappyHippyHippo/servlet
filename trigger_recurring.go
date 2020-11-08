package servlet

import (
	"fmt"
	"time"
)

// TriggerRecurring defines a trigger instance used to execute a process
// periodically with a defined frequency.
type TriggerRecurring struct {
	Trigger
}

// NewTriggerRecurring instantiate a new trigger that will execute a
// callback method recurrently with a defined periodicity.
func NewTriggerRecurring(period time.Duration, callback TriggerCallback) (*TriggerRecurring, error) {
	if callback == nil {
		return nil, fmt.Errorf("invalid nil 'callback' argument")
	}

	t := &TriggerRecurring{
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
