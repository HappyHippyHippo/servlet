package servlet

import (
	"testing"
)

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
