package trigger

import (
	"fmt"
	"time"
)

type recurringTrigger struct {
	trigger
}

// NewRecurringTrigger intantiate a new trigger that will execute a
// callback method recurrently with a defined periodicity.
func NewRecurringTrigger(period time.Duration, callback Callback) (Trigger, error) {
	if callback == nil {
		return nil, fmt.Errorf("Invalid nil 'callback' argument")
	}

	t := &recurringTrigger{
		trigger: trigger{
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

// Close method will stop the trigger from execution.
func (t *recurringTrigger) Close() error {
	return t.Stop()
}
