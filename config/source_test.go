package config

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_Source_Close(t *testing.T) {
	s := &source{&sync.RWMutex{}, partial{}}

	t.Run("no-op", func(t *testing.T) {
		if err := s.Close(); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Source_Has(t *testing.T) {
	search := "path"
	expected := true

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mutex := NewMockRWMutex(ctrl)
	mutex.EXPECT().RLock().Times(1)
	mutex.EXPECT().RUnlock().Times(1)

	partial := NewMockPartial(ctrl)
	partial.EXPECT().Has(search).Return(expected).Times(1)

	s := &source{mutex, partial}

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		if value := s.Has(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}

func Test_Source_Get(t *testing.T) {
	search := "path"
	expected := "value"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mutex := NewMockRWMutex(ctrl)
	mutex.EXPECT().RLock().Times(1)
	mutex.EXPECT().RUnlock().Times(1)

	partial := NewMockPartial(ctrl)
	partial.EXPECT().Get(search).Return(expected).Times(1)

	s := &source{mutex, partial}

	t.Run("lock and redirect to the stored partial", func(t *testing.T) {
		if value := s.Get(search); value != expected {
			t.Errorf("returned the (%v) value", value)
		}
	})
}
