package trigger

import (
	"testing"
	"time"
)

func Test_NewPulseTrigger(t *testing.T) {
	t.Run("should return nil when missing callback", func(t *testing.T) {
		action := "Creating a new pulse trigger without callback function"

		expected := "Invalid nil 'callback' argument"

		trigger, err := NewPulseTrigger(20*time.Millisecond, nil)

		if trigger != nil {
			trigger.Stop()
			t.Errorf("%s returned an unexpected valid reference to a trigger, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s returned the error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("creates a new pulse trigger", func(t *testing.T) {
		action := "Creating a new pulse trigger"

		trigger, err := NewPulseTrigger(20*time.Millisecond, func() error { return nil })

		if trigger == nil {
			t.Errorf("%s didn't return a valid reference to a new pulse trigger", action)
		} else {
			trigger.Stop()
		}
		if err != nil {
			t.Errorf("%s returned an unexpected error : %v", action, err)
		}
	})
}

func Test_PulseTrigger_Close(t *testing.T) {
	t.Run("is the same as stopping it", func(t *testing.T) {
		action := "Closing the trigger prior to the wait time duration"

		check := false

		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = true; return nil })
		trigger.Close()

		time.Sleep(5 * time.Millisecond)

		if check {
			t.Errorf("%s didn't prevent the trigger execution", action)
		}
	})
}

func Test_PulseTrigger_Duration(t *testing.T) {
	t.Run("retrieves the initialization wait duration", func(t *testing.T) {
		action := "Retrieving the wait duration"

		duration := 20 * time.Millisecond

		trigger, _ := NewPulseTrigger(duration, func() error { return nil })
		trigger.Stop()

		if result := trigger.Timer(); result != duration {
			t.Errorf("%s returned (%v) wait duration, expected (%v)", action, result, duration)
		}
	})
}

func Test_PulseTrigger_IsStopped(t *testing.T) {
	t.Run("should return false if called after creation", func(t *testing.T) {
		action := "Retrieving the stopped state after creation"

		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Stop()

		if trigger.IsStopped() {
			t.Errorf("%s returned true, expected a false value", action)
		}
	})

	t.Run("should return true after calling Stop method", func(t *testing.T) {
		action := "Retrieving the stopped state after stopping"

		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { return nil })
		trigger.Stop()

		if !trigger.IsStopped() {
			t.Errorf("%s returned false, expected a true value", action)
		}
	})
}

func Test_PulseTrigger_Stop(t *testing.T) {
	t.Run("should prevent triggering if called prior to the first execution", func(t *testing.T) {
		action := "Stopping the trigger prior to the wait time duration"

		check := false

		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = true; return nil })
		trigger.Stop()

		time.Sleep(50 * time.Millisecond)

		if check {
			t.Errorf("%s didn't prevent the trigger execution", action)
		}
	})
}

func Test_PulseTrigger(t *testing.T) {
	t.Run("should only trigger execution once", func(t *testing.T) {
		action := "Waiting for trigger after wait time duration"

		check := 0

		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = check + 1; return nil })
		defer trigger.Stop()

		time.Sleep(100 * time.Millisecond)

		if check == 0 {
			t.Errorf("%s didn't called the callback function once", action)
		}
		if check > 1 {
			t.Errorf("%s recurrently called the callback function, expected only once", action)
		}
	})
}
