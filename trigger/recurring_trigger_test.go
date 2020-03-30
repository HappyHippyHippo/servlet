package trigger

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewRecurringTrigger(t *testing.T) {
	t.Run("should return nil when missing callback", func(t *testing.T) {
		action := "Creating a new recurring trigger without a callback function"

		expected := "Invalid nil 'callback' argument"

		trigger, err := NewRecurringTrigger(20*time.Millisecond, nil)

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

	t.Run("creats a new recurring trigger", func(t *testing.T) {
		action := "Creating a new recurring trigger"

		trigger, err := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })

		if trigger == nil {
			t.Errorf("%s didn't return a valid reference to a new recurring trigger", action)
		} else {
			trigger.Stop()
		}
		if err != nil {
			t.Errorf("%s returned an unexpected error : %v", action, err)
		}
	})
}

func Test_RecurringTrigger_Close(t *testing.T) {
	t.Run("is the same as stopping it", func(t *testing.T) {
		action := "Closing the trigger prior to the interval time duration"

		check := false

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = true; return nil })
		trigger.Close()

		if check {
			t.Errorf("%s didn't stop the trigger to be executed", action)
		}
	})
}

func Test_RecurringTrigger_Duration(t *testing.T) {
	t.Run("retrieves a initialization interval duration", func(t *testing.T) {
		action := "Retrieving a recursive trigger interval duration"

		duration := 20 * time.Millisecond

		trigger, _ := NewRecurringTrigger(duration, func() error { return nil })
		defer trigger.Stop()

		if result := trigger.Timer(); result != duration {
			t.Errorf("%s returned (%v) interval duration, expected (%v)", action, result, duration)
		}
	})
}

func Test_RecurringTrigger_IsStopped(t *testing.T) {
	t.Run("should return false if called after creation", func(t *testing.T) {
		action := "Retrieving a recursive trigger stopped state after creation"

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Stop()

		if trigger.IsStopped() {
			t.Errorf("%s returned true, expected a false value", action)
		}
	})

	t.Run("should return true after calling Stop method", func(t *testing.T) {
		action := "Retrieving a recursive trigger stopped state after stopping"

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })
		trigger.Stop()

		if !trigger.IsStopped() {
			t.Errorf("%s returned false, expected a true value", action)
		}
	})
}

func Test_RecurringTrigger_Stop(t *testing.T) {
	t.Run("should prevent triggering if called prior to first execution", func(t *testing.T) {
		action := "Stopping the trigger prior to the interval time duration"

		check := false

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = true; return nil })
		trigger.Stop()

		if check {
			t.Errorf("%s didn't stop the trigger to be executed", action)
		}
	})
}

func Test_RecurringTrigger(t *testing.T) {
	t.Run("should run trigger multiple times", func(t *testing.T) {
		action := "Waiting for trigger after interval time duration"

		check := 0

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = check + 1; return nil })
		defer trigger.Stop()

		time.Sleep(100 * time.Millisecond)

		if check <= 2 {
			t.Errorf("%s didn't recurrently called the callback function", action)
		}
	})

	t.Run("should stop the trigger on callback error", func(t *testing.T) {
		action := "Having a callback function return an error"

		check := 0

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = check + 1; return fmt.Errorf("__dummy_error__") })
		defer trigger.Stop()

		time.Sleep(100 * time.Millisecond)

		if check != 1 {
			t.Errorf("%s didn't stop recursion after the first error", action)
		}
	})
}
