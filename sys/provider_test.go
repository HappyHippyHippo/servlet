package sys

import (
	"reflect"
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

func Test_NewProvider(t *testing.T) {
	parameters := NewDefaultParameters()
	t.Run("creates a new provider", func(t *testing.T) {
		if provider := NewProvider(parameters).(*provider); provider == nil {
			t.Errorf("didn't return a valid reference")
		} else if !reflect.DeepEqual(parameters, provider.params) {
			t.Errorf("stored (%v) parameters", provider.params)
		}
	})
}

func Test_Provider_Register(t *testing.T) {
	matcher := &factoryMatcher{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	provider := NewProvider(NewDefaultParameters())

	container := NewMockContainer(ctrl)
	container.EXPECT().Add(ContainerFileSystemID, matcher)

	t.Run("register the file system adapter in the container", func(t *testing.T) {
		provider.Register(container)

		entity := matcher.factory(container)
		switch entity.(type) {
		case afero.Fs:
		default:
			t.Errorf("didn't return a file system adapter")
		}
	})
}
