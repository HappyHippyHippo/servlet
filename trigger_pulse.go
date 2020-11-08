package servlet

import (
	"fmt"
	"time"
)

// TriggerPulse defines a trigger instance used to execute a process once
// after a defines period of time.
type TriggerPulse struct {
	Trigger
}

// NewTriggerPulse instantiate a new pulse trigger that will execute a
// callback method after a determined amount of time.
func NewTriggerPulse(delay time.Duration, callback TriggerCallback) (*TriggerPulse, error) {
	if callback == nil {
		return nil, fmt.Errorf("invalid nil 'callback' argument")
	}

	t := &TriggerPulse{
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
