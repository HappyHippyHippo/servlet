package log

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/servlet"
	"github.com/happyhippyhippo/servlet/config"
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
	t.Run("should creates a new log provider with default parameters if none are given", func(t *testing.T) {
		action := "Creating a new log provider without parameters"

		provider := NewProvider(nil).(*provider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new log provider", action)
		}

		if check := provider.params.id; check != ContainerID {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, ContainerID, check)
		}
		if check := provider.params.fileSystemID; check != sys.ContainerFileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, sys.ContainerFileSystemID, check)
		}
		if check := provider.params.configID; check != config.ContainerID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, config.ContainerID, check)
		}
		if check := provider.params.formatterFactoryID; check != ContainerFormatterFactoryID {
			t.Errorf("%s didn't store the expected (%v) formatter factory id, returned (%v)", action, ContainerFormatterFactoryID, check)
		}
		if check := provider.params.streamFactoryID; check != ContainerStreamFactoryID {
			t.Errorf("%s didn't store the expected (%v) stream factory id, returned (%v)", action, ContainerStreamFactoryID, check)
		}
		if check := provider.params.loaderID; check != ContainerLoaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, ContainerLoaderID, check)
		}
	})

	t.Run("should creates a new log provider with given parameters", func(t *testing.T) {
		action := "Creating a new log provider with parameters"

		id := "__dummy_id__"
		fileSystemID := "__dummy_file_system_id__"
		configID := "__dummy_config_id__"
		formatterFactoryID := "__dummy_formatter_factory_id__"
		streamFactoryID := "__dummy_stream_factory_id__"
		loaderID := "__dummy_loader_id__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		parameters := NewMockProviderParameters(ctrl)
		parameters.EXPECT().GetID().Return(id).Times(1)
		parameters.EXPECT().GetFileSystemID().Return(fileSystemID).Times(1)
		parameters.EXPECT().GetConfigID().Return(configID).Times(1)
		parameters.EXPECT().GetFormatterFactoryID().Return(formatterFactoryID).Times(1)
		parameters.EXPECT().GetStreamFactoryID().Return(streamFactoryID).Times(1)
		parameters.EXPECT().GetLoaderID().Return(loaderID).Times(1)

		provider := NewProvider(parameters).(*provider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new log provider", action)
		}
		if check := provider.params.id; check != id {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, id, check)
		}
		if check := provider.params.fileSystemID; check != fileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, fileSystemID, check)
		}
		if check := provider.params.configID; check != configID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, configID, check)
		}
		if check := provider.params.formatterFactoryID; check != formatterFactoryID {
			t.Errorf("%s didn't store the expected (%v) formatter factory id, returned (%v)", action, formatterFactoryID, check)
		}
		if check := provider.params.streamFactoryID; check != streamFactoryID {
			t.Errorf("%s didn't store the expected (%v) stream factory id, returned (%v)", action, streamFactoryID, check)
		}
		if check := provider.params.loaderID; check != loaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, loaderID, check)
		}
	})
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should retrieve the formatter factory", func(t *testing.T) {
		action := "calling the formatter factory builded by the registed provider"

		formatterFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, formatterFactoryMatcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := formatterFactoryMatcher.factory(container)
		switch entity.(type) {
		case FormatterFactory:
		default:
			t.Errorf("%s didn't return a valid reference to a new formatter factory", action)
		}
	})

	t.Run("should panic when trying to retrieve a stream factory when the file system adapter is missing", func(t *testing.T) {
		action := "calling the stream factory builded by the registed provider when the file system adapter is missing"

		streamFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, streamFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		streamFactoryMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a stream factory when the formatter factory is missing", func(t *testing.T) {
		action := "calling the stream factory builded by the registed provider when the formatter factory is missing"

		streamFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, streamFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		streamFactoryMatcher.factory(container)
	})

	t.Run("should retrieve the stream factory", func(t *testing.T) {
		action := "calling the stream factory builded by the registed provider"

		streamFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		formatterFactory := NewMockFormatterFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, streamFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := streamFactoryMatcher.factory(container)
		switch entity.(type) {
		case StreamFactory:
		default:
			t.Errorf("%s didn't return a valid reference to a new stream factory", action)
		}
	})

	t.Run("should retrieve the logger", func(t *testing.T) {
		action := "calling the logger builded by the registed provider"

		loggerMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, loggerMatcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := loggerMatcher.factory(container)
		switch entity.(type) {
		case Logger:
		default:
			t.Errorf("%s didn't return a valid reference to a new logger", action)
		}
	})

	t.Run("should panic when trying to retrieve a loader when the formatter factory is missing", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the formatter factory is missing"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the formatter factory is non-compliant with the interface", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the formatter factory is non-compliant with the interface"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the stream factory is missing", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the stream factory is missing"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the stream factory is non-compliant with the interface", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the stream factory is non-compliant with the interface"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
			container.EXPECT().Get(ContainerStreamFactoryID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the logger is missing", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the logger is missing"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(streamFactory),
			container.EXPECT().Get(ContainerID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the logger is non-compliant with the interface", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the logger is non-compliant with the interface"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(streamFactory),
			container.EXPECT().Get(ContainerID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerStreamFactoryID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should retrieve the loader", func(t *testing.T) {
		action := "calling the loader builded by the registed provider"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerFormatterFactoryID, matcher),
			container.EXPECT().Add(ContainerStreamFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerFormatterFactoryID).Return(formatterFactory),
			container.EXPECT().Get(ContainerStreamFactoryID).Return(streamFactory),
			container.EXPECT().Get(ContainerID).Return(logger),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := loaderMatcher.factory(container)
		switch entity.(type) {
		case Loader:
		default:
			t.Errorf("%s didn't return a valid reference to a new logger", action)
		}
	})
}

func Test_Provider_Boot(t *testing.T) {
	t.Run("should panic if the loader is not registered", func(t *testing.T) {
		action := "calling the loader builded by the registed provider"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(nil)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerLoaderID)
			}
		}()

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})

	t.Run("should panic if the loader is non-compliant with the interface", func(t *testing.T) {
		action := "calling a non-compliant loader builded by the registed provider"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return("__something_else__")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerLoaderID)
			}
		}()

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})

	t.Run("should panic if the config is not registered", func(t *testing.T) {
		action := "calling the config builded by the registed provider"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(loader)
		container.EXPECT().Get(config.ContainerID).Return(nil)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerLoaderID)
			}
		}()

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})

	t.Run("should panic if the config is non-compliant with the interface", func(t *testing.T) {
		action := "calling a non-compliant config builded by the registed provider"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(loader)
		container.EXPECT().Get(config.ContainerID).Return("__something_else__")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerLoaderID)
			}
		}()

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})

	t.Run("should load the configuration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cfg := NewMockConfig(ctrl)
		loader := NewMockLoader(ctrl)
		loader.EXPECT().Load(cfg).Times(1)

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(loader)
		container.EXPECT().Get(config.ContainerID).Return(cfg)

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})
}
