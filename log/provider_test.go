package log

import (
	reflect "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet"
	config "github.com/happyhippyhippo/servlet/config"
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

	t.Run("register the formatter factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case FormatterFactory:
		default:
			t.Errorf("didn't return a formatter factory")
		}
	})

	t.Run("panic retrieving stream factory with missing file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
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

	t.Run("panic retrieving stream factory with invalid file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
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

	t.Run("panic retrieving stream factory with missing formatter factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(nil),
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

	t.Run("panic retrieving stream factory with invalid formatter factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerFormatterFactoryID).Return("invalid_reference"),
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

	t.Run("register the stream factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, checkMatcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case *streamFactory:
		default:
			t.Errorf("didn't return a stream factory")
		}
	})

	t.Run("register the logger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, checkMatcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(parameters)
		provider.Register(container)

		entity := checkMatcher.factory(container)
		switch entity.(type) {
		case *logger:
		default:
			t.Errorf("didn't return a logger")
		}
	})

	t.Run("panic retrieving loader with missing logger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		container.EXPECT().Get(ContainerLoggerID).Return(nil)

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving loader with invalid logger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		container.EXPECT().Get(ContainerLoggerID).Return("invalid_reference")

		provider := NewProvider(parameters)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		checkMatcher.factory(container)
	})

	t.Run("panic retrieving loader with missing stream factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewMockLogger(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoggerID).Return(logger),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(nil),
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

	t.Run("panic retrieving loader with invalid stream factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewMockLogger(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoggerID).Return(logger),
			container.EXPECT().Get(ContainerStreamFactoryID).Return("invalid_reference"),
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

	t.Run("register the loader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewMockLogger(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerLoggerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, checkMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoggerID).Return(logger),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(streamFactory),
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

	t.Run("panic when missing logger", func(t *testing.T) {
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

	t.Run("panic when loader is invalid", func(t *testing.T) {
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

	t.Run("panic when missing config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoaderID).Return(loader),
			container.EXPECT().Get(config.ContainerConfigID).Return(nil),
		)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		provider := NewProvider(parameters)
		provider.Boot(container)
	})

	t.Run("panic when config is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoaderID).Return(loader),
			container.EXPECT().Get(config.ContainerConfigID).Return("invalid_reference"),
		)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("did not panic")
			}
		}()

		provider := NewProvider(parameters)
		provider.Boot(container)
	})

	t.Run("load configuration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockConfig(ctrl)
		loader := NewMockLoader(ctrl)
		loader.EXPECT().Load(conf).Times(1)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Get(ContainerLoaderID).Return(loader),
			container.EXPECT().Get(config.ContainerConfigID).Return(conf),
		)

		provider := NewProvider(parameters)
		provider.Boot(container)
	})
}
