package trigger

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewRecurringTrigger(t *testing.T) {
	t.Run("error when missing callback", func(t *testing.T) {
		if trigger, err := NewRecurringTrigger(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "Invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new recurring trigger", func(t *testing.T) {
		if trigger, err := NewRecurringTrigger(20*time.Millisecond, func() error { return nil }); trigger == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			defer trigger.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_RecurringTrigger_Close(t *testing.T) {
	check := false
	trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = true; return nil })
	trigger.Close()

	t.Run("is the same as stopping it", func(t *testing.T) {
		time.Sleep(40 * time.Millisecond)
		if check {
			t.Errorf("didn't stop the trigger to be executed")
		}
	})
}

func Test_RecurringTrigger_Timer(t *testing.T) {
	duration := 20 * time.Millisecond
	trigger, _ := NewRecurringTrigger(duration, func() error { return nil })
	defer trigger.Close()

	t.Run("retrieves the trigger interval duration", func(t *testing.T) {
		if result := trigger.Timer(); result != duration {
			t.Errorf("returned (%v) interval duration", result)
		}
	})
}

func Test_RecurringTrigger_IsStopped(t *testing.T) {
	t.Run("return false if called after creation", func(t *testing.T) {
		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		if trigger.IsStopped() {
			t.Errorf("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		trigger.Stop()

		if !trigger.IsStopped() {
			t.Errorf("returned false")
		}
	})
}

func Test_RecurringTrigger_Stop(t *testing.T) {
	check := false

	trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = true; return nil })
	defer trigger.Close()
	trigger.Stop()

	t.Run("prevent triggering if called prior to first execution", func(t *testing.T) {
		time.Sleep(40 * time.Millisecond)
		if check {
			t.Errorf("didn't stop the trigger to be executed")
		}
	})
}

func Test_RecurringTrigger(t *testing.T) {
	t.Run("run trigger multiple times", func(t *testing.T) {
		check := 0

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = check + 1; return nil })
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check <= 2 {
			t.Errorf("didn't recurrently called the callback function")
		}
	})

	t.Run("stop the trigger on callback error", func(t *testing.T) {
		check := 0

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = check + 1; return fmt.Errorf("__dummy_error__") })
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check != 1 {
			t.Errorf("didn't stop recursion calls after the first error")
		}
	})
}
