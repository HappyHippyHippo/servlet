package servlet

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewTriggerRecurring(t *testing.T) {
	t.Run("nil callback", func(t *testing.T) {
		if trigger, err := NewTriggerRecurring(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new recurring trigger", func(t *testing.T) {
		if trigger, err := NewTriggerRecurring(20*time.Millisecond, func() error {
			return nil
		}); trigger == nil {
			t.Error("didn't returned a valid reference")
		} else {
			defer trigger.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_TriggerRecurring_Close(t *testing.T) {
	t.Run("is the same as stopping it", func(t *testing.T) {
		check := false
		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error {
			check = true
			return nil
		})
		trigger.Close()

		time.Sleep(40 * time.Millisecond)
		if check {
			t.Error("didn't stop the trigger to be executed")
		}
	})
}

func Test_TriggerRecurring_Timer(t *testing.T) {
	t.Run("retrieves the trigger interval duration", func(t *testing.T) {
		duration := 20 * time.Millisecond
		trigger, _ := NewTriggerRecurring(duration, func() error {
			return nil
		})
		defer trigger.Close()

		if result := trigger.Timer(); result != duration {
			t.Errorf("returned (%v) interval duration", result)
		}
	})
}

func Test_TriggerRecurring_IsStopped(t *testing.T) {
	t.Run("return false if called after creation", func(t *testing.T) {
		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error {
			return nil
		})
		defer trigger.Close()

		if trigger.IsStopped() {
			t.Error("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error {
			return nil
		})
		defer trigger.Close()

		trigger.Stop()

		if !trigger.IsStopped() {
			t.Error("returned false")
		}
	})
}

func Test_TriggerRecurring_Stop(t *testing.T) {
	t.Run("prevent triggering if called prior to first execution", func(t *testing.T) {
		check := false

		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error { check = true; return nil })
		defer trigger.Close()
		trigger.Stop()

		time.Sleep(40 * time.Millisecond)
		if check {
			t.Error("didn't stop the trigger to be executed")
		}
	})
}

func Test_TriggerRecurring(t *testing.T) {
	t.Run("run trigger multiple times", func(t *testing.T) {
		check := 0

		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error { check = check + 1; return nil })
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check <= 2 {
			t.Error("didn't recurrently called the callback function")
		}
	})

	t.Run("stop the trigger on callback error", func(t *testing.T) {
		check := 0

		trigger, _ := NewTriggerRecurring(20*time.Millisecond, func() error { check = check + 1; return fmt.Errorf("__dummy_error__") })
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check != 1 {
			t.Error("didn't stop recursion calls after the first error")
		}
	})
}
