package servlet

import (
	"testing"
	"time"
)

func Test_NewTriggerPulse(t *testing.T) {
	t.Run("nil callback", func(t *testing.T) {
		if trigger, err := NewTriggerPulse(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new pulse trigger", func(t *testing.T) {
		if trigger, err := NewTriggerPulse(20*time.Millisecond, func() error {
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

func Test_TriggerPulse_Close(t *testing.T) {
	t.Run("is the same as stopping it", func(t *testing.T) {
		check := false

		trigger, _ := NewTriggerPulse(20*time.Millisecond, func() error {
			check = true
			return nil
		})
		trigger.Close()

		time.Sleep(40 * time.Millisecond)
		if check {
			t.Error("didn't prevented the trigger to execute")
		}
	})
}

func Test_TriggerPulse_Timer(t *testing.T) {
	t.Run("retrieves the trigger time", func(t *testing.T) {
		duration := 20 * time.Millisecond

		trigger, _ := NewTriggerPulse(duration, func() error {
			return nil
		})
		defer trigger.Close()

		if result := trigger.Timer(); result != duration {
			t.Errorf("returned (%v) wait duration", result)
		}
	})
}

func Test_TriggerPulse_IsStopped(t *testing.T) {
	t.Run("return false if called after creation", func(t *testing.T) {
		trigger, _ := NewTriggerPulse(20*time.Millisecond, func() error {
			return nil
		})
		defer trigger.Close()

		if trigger.IsStopped() {
			t.Error("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewTriggerPulse(20*time.Millisecond, func() error {
			return nil
		})
		defer trigger.Close()

		trigger.Stop()
		if !trigger.IsStopped() {
			t.Error("returned false")
		}
	})
}

func Test_TriggerPulse_Stop(t *testing.T) {
	t.Run("prevent triggering if called prior to the first execution", func(t *testing.T) {
		check := false

		trigger, _ := NewTriggerPulse(20*time.Millisecond, func() error {
			check = true
			return nil
		})
		defer trigger.Close()
		trigger.Stop()

		time.Sleep(40 * time.Millisecond)
		if check {
			t.Error("didn't prevented the trigger to execute")
		}
	})
}

func Test_TriggerPulse(t *testing.T) {
	t.Run("only trigger execution once", func(t *testing.T) {
		check := 0

		trigger, _ := NewTriggerPulse(20*time.Millisecond, func() error {
			check = check + 1
			return nil
		})
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check == 0 {
			t.Error("didn't called the callback function once")
		} else if check > 1 {
			t.Error("recurrently called the callback function")
		}
	})
}
