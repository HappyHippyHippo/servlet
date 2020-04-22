package config

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/sys"
)

type containerFactoryMatcher struct {
	factory servlet.Factory
}

func (m *containerFactoryMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case servlet.Factory:
		m.factory = x.(servlet.Factory)
		return true
	}
	return false
}

func (m *containerFactoryMatcher) String() string {
	return "a container factory"
}

func Test_NewProvider(t *testing.T) {
	t.Run("creates a new provider", func(t *testing.T) {
		parameters := NewDefaultParameters()
		if provider := NewProvider(parameters).(*provider); provider == nil {
			t.Errorf("didn't return a valid reference")
		} else if !reflect.DeepEqual(parameters, provider.params) {
			t.Errorf("stored (%v) parameters", provider.params)
		}
	})
}

func Test_Provider_Register(t *testing.T) {
	parameters := NewDefaultParameters()
	checkMatcher := &containerFactoryMatcher{}
	matcher := &containerFactoryMatcher{}

	t.Run("register the decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case DecoderFactory:
		default:
			t.Errorf("didn't return a decoder factory")
		}
	})

	t.Run("panic retrieving source factory with missing file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		container.EXPECT().Get(sys.ContainerFileSystemID).Return(nil)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving source factory with invalid file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		container.EXPECT().Get(sys.ContainerFileSystemID).Return("invalid_reference")

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving source factory with missing decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return(nil),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving source factory with invalid decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return("invalid_reference"),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("register the source factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return(decoderFactory),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case *sourceFactory:
		default:
			t.Errorf("didn't return a source factory")
		}
	})

	t.Run("register the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, checkMatcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case *config:
		default:
			t.Errorf("didn't return a config")
		}
	})

	t.Run("panic retrieving loader with missing config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerConfigID).Return(nil),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving loader with invalid config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerConfigID).Return("invalid_reference"),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving loader with missing source factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerConfigID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return(nil),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving loader with invalid source factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerConfigID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return("invalid_reference"),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("retrieving loader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerConfigID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerConfigID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return(sourceFactory),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case *loader:
		default:
			t.Errorf("didn't return a loader")
		}
	})
}

func Test_Provider_Boot(t *testing.T) {
	parameters := NewDefaultParameters()

	t.Run("panic if the loader is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(nil)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		provider := NewProvider(parameters)
		provider.Boot(container)
	})

	t.Run("panic if the loader is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return("invalid_reference")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		provider := NewProvider(parameters)
		provider.Boot(container)
	})

	t.Run("load the configuration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)
		loader.EXPECT().Load(BaseSourceID, BaseSourcePath, BaseSourceFormat).Times(1)

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(loader)

		provider := NewProvider(parameters)
		provider.Boot(container)
	})
}
