package trigger

import (
	"fmt"
	"time"
)

type pulseTrigger struct {
	trigger
}

// NewPulseTrigger intantiate a new pulse trigger that will execute a
// callback method after a determined amount of time.
func NewPulseTrigger(delay time.Duration, callback Callback) (Trigger, error) {
	if callback == nil {
		return nil, fmt.Errorf("Invalid nil 'callback' argument")
	}

	t := &pulseTrigger{
		trigger: trigger{
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
					t.callback()
				}
				return
			case <-t.channelStop:
				return
			}
		}
	}()

	return t, nil
}

// Close method will stop the trigger, this means that if the trigger have not
// yet been executed, then the associated process will not execute.
func (t *pulseTrigger) Close() error {
	return t.Stop()
}
