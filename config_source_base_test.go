package servlet

import (
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func Test_ConfigSourceBase_Close(t *testing.T) {
	s := &ConfigSourceBase{&sync.Mutex{}, ConfigPartial{}}
	s.Close()
}

func Test_ConfigSourceBase_Has(t *testing.T) {
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

		var source *ConfigSourceBase
		_ = source.Has("path")
	})

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		search := "path"
		expected := true
		partial := ConfigPartial{search: "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mutex := NewMockLocker(ctrl)
		mutex.EXPECT().Lock().Times(1)
		mutex.EXPECT().Unlock().Times(1)

		s := &ConfigSourceBase{mutex: mutex, partial: partial}

		if value := s.Has(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}

func Test_ConfigSourceBase_Get(t *testing.T) {
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

		var source *ConfigSourceBase
		_ = source.Get("path")
	})

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		search := "path"
		expected := "value"
		partial := ConfigPartial{search: expected}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mutex := NewMockLocker(ctrl)
		mutex.EXPECT().Lock().Times(1)
		mutex.EXPECT().Unlock().Times(1)

		s := &ConfigSourceBase{mutex: mutex, partial: partial}

		if value := s.Get(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}
