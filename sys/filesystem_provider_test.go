package sys

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet"
	"github.com/spf13/afero"
)

type factoryMatcher struct {
	factory servlet.Factory
}

func (m *factoryMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case servlet.Factory:
		m.factory = x.(servlet.Factory)
		return true
	}
	return false
}

func (m *factoryMatcher) String() string {
	return "a container factory"
}

func Test_NewFileSystemProvider(t *testing.T) {
	t.Run("should creates a new provider with default parameters if none are given", func(t *testing.T) {
		action := "Creating a new provider without parameters"

		provider := NewFileSystemProvider(nil).(*fileSystemProvider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new provider", action)
		}

		if check := provider.id; check != "servlet.filesystem" {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, ContainerFileSystemID, check)
		}
	})

	t.Run("should creates a new provider with given parameters", func(t *testing.T) {
		action := "Creating a new provider with parameters"

		id := "__dummy_id__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		parameters := NewMockFileSystemProviderParameters(ctrl)
		parameters.EXPECT().GetID().Return(id).Times(1)

		provider := NewFileSystemProvider(parameters).(*fileSystemProvider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new file system provider", action)
		}
		if check := provider.id; check != id {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, id, check)
		}
	})
}

func Test_FileSystemProvider_Register(t *testing.T) {
	t.Run("should register a container factory that returns a file system adapter", func(t *testing.T) {
		action := "calling the entity builded by the registed provider"

		matcher := &factoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		container.EXPECT().Add(ContainerFileSystemID, matcher)

		provider := NewFileSystemProvider(nil)
		provider.Register(container)

		entity := matcher.factory(container)
		switch entity.(type) {
		case afero.Fs:
		default:
			t.Errorf("%s didn't return a valid reference to a new provider", action)
		}
	})
}
