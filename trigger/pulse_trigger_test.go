package trigger

import (
	"testing"
	"time"
)

func Test_NewPulseTrigger(t *testing.T) {
	t.Run("error when missing callback", func(t *testing.T) {
		if trigger, err := NewPulseTrigger(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new pulse trigger", func(t *testing.T) {
		if trigger, err := NewPulseTrigger(20*time.Millisecond, func() error { return nil }); trigger == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			defer trigger.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})
}

func Test_PulseTrigger_Close(t *testing.T) {
	check := false

	trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = true; return nil })
	trigger.Close()

	t.Run("is the same as stopping it", func(t *testing.T) {
		time.Sleep(40 * time.Millisecond)
		if check {
			t.Errorf("didn't prevent the trigger to execute")
		}
	})
}

func Test_PulseTrigger_Timer(t *testing.T) {
	duration := 20 * time.Millisecond

	trigger, _ := NewPulseTrigger(duration, func() error { return nil })
	defer trigger.Close()

	t.Run("retrieves the trigger time", func(t *testing.T) {
		if result := trigger.Timer(); result != duration {
			t.Errorf("returned (%v) wait duration", result)
		}
	})
}

func Test_PulseTrigger_IsStopped(t *testing.T) {
	t.Run("return false if called after creation", func(t *testing.T) {
		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		if trigger.IsStopped() {
			t.Errorf("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		trigger.Stop()
		if !trigger.IsStopped() {
			t.Errorf("returned false")
		}
	})
}

func Test_PulseTrigger_Stop(t *testing.T) {
	check := false

	trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = true; return nil })
	defer trigger.Close()
	trigger.Stop()

	t.Run("prevent triggering if called prior to the first execution", func(t *testing.T) {
		time.Sleep(40 * time.Millisecond)
		if check {
			t.Errorf("didn't prevent the trigger to execute")
		}
	})
}

func Test_PulseTrigger(t *testing.T) {
	check := 0

	trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { check = check + 1; return nil })
	defer trigger.Close()

	t.Run("only trigger execution once", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)

		if check == 0 {
			t.Errorf("didn't called the callback function once")
		} else if check > 1 {
			t.Errorf("recurrently called the callback function")
		}
	})
}
