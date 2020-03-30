package config

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_Source_Close(t *testing.T) {
	t.Run("should lock the source and redirect the request to the stored partial", func(t *testing.T) {
		action := "Calling the close method over the base config source"

		s := &source{&sync.RWMutex{}, partial{}}

		if check := s.Close(); check != nil {
			t.Errorf("%s returned the unexpected error : %v", action, check)
		}
	})
}

func Test_Source_Has(t *testing.T) {
	t.Run("should lock the source and redirect the request to the stored partial", func(t *testing.T) {
		action := "Checking if a path is present in the source"

		path := "__dummy_path__"
		expected := true

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().Has(path).Return(expected).Times(1)
		mutex := NewMockRWMutex(ctrl)
		mutex.EXPECT().RLock().Times(1)
		mutex.EXPECT().RUnlock().Times(1)

		s := &source{mutex, partial}

		if check := s.Has(path); check != expected {
			t.Errorf("%s return %v, expected %v", action, check, expected)
		}
	})
}

func Test_Source_Get(t *testing.T) {
	t.Run("should lock the source and redirect the request to the stored partial", func(t *testing.T) {
		action := "Checking if a path is present in the source"

		path := "__dummy_path__"
		expected := "__dummy_value__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().Get(path).Return(expected).Times(1)
		mutex := NewMockRWMutex(ctrl)
		mutex.EXPECT().RLock().Times(1)
		mutex.EXPECT().RUnlock().Times(1)

		s := &source{mutex, partial}

		if check := s.Get(path); check != expected {
			t.Errorf("%s return %v, expected %v", action, check, expected)
		}
	})
}
