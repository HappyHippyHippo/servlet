package servlet

import (
	"fmt"
	"testing"
	"time"
)

/// ---------------------------------------------------------------------------
/// Trigger
/// ---------------------------------------------------------------------------

func Test_Trigger_Close(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var trigger *Trigger
		trigger.Close()
	})
}

func Test_Trigger_Stop(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var trigger *Trigger
		trigger.Stop()
	})
}

/// ---------------------------------------------------------------------------
/// PulseTrigger
/// ---------------------------------------------------------------------------

func Test_NewPulseTrigger(t *testing.T) {
	t.Run("nil callback", func(t *testing.T) {
		if trigger, err := NewPulseTrigger(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new pulse trigger", func(t *testing.T) {
		if trigger, err := NewPulseTrigger(20*time.Millisecond, func() error { return nil }); trigger == nil {
			t.Error("didn't returned a valid reference")
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
			t.Error("didn't prevented the trigger to execute")
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
			t.Error("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewPulseTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		trigger.Stop()
		if !trigger.IsStopped() {
			t.Error("returned false")
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
			t.Error("didn't prevented the trigger to execute")
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
			t.Error("didn't called the callback function once")
		} else if check > 1 {
			t.Error("recurrently called the callback function")
		}
	})
}

/// ---------------------------------------------------------------------------
/// RecurringTrigger
/// ---------------------------------------------------------------------------

func Test_NewRecurringTrigger(t *testing.T) {
	t.Run("nil callback", func(t *testing.T) {
		if trigger, err := NewRecurringTrigger(20*time.Millisecond, nil); trigger != nil {
			defer trigger.Close()
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'callback' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new recurring trigger", func(t *testing.T) {
		if trigger, err := NewRecurringTrigger(20*time.Millisecond, func() error { return nil }); trigger == nil {
			t.Error("didn't returned a valid reference")
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
			t.Error("didn't stop the trigger to be executed")
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
			t.Error("returned true")
		}
	})

	t.Run("return true after calling Stop method", func(t *testing.T) {
		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { return nil })
		defer trigger.Close()

		trigger.Stop()

		if !trigger.IsStopped() {
			t.Error("returned false")
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
			t.Error("didn't stop the trigger to be executed")
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
			t.Error("didn't recurrently called the callback function")
		}
	})

	t.Run("stop the trigger on callback error", func(t *testing.T) {
		check := 0

		trigger, _ := NewRecurringTrigger(20*time.Millisecond, func() error { check = check + 1; return fmt.Errorf("__dummy_error__") })
		defer trigger.Close()

		time.Sleep(100 * time.Millisecond)

		if check != 1 {
			t.Error("didn't stop recursion calls after the first error")
		}
	})
}
