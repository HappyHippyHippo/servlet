package config

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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
	t.Run("should creates a new config provider with default parameters if none are given", func(t *testing.T) {
		action := "Creating a new config provider without parameters"

		provider := NewProvider(nil).(*provider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new config provider", action)
		}

		if check := provider.params.id; check != ContainerID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, ContainerID, check)
		}
		if check := provider.params.fileSystemID; check != sys.ContainerFileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, sys.ContainerFileSystemID, check)
		}
		if check := provider.params.loaderID; check != ContainerLoaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, ContainerLoaderID, check)
		}
		if check := provider.params.sourceFactoryID; check != ContainerSourceFactoryID {
			t.Errorf("%s didn't store the expected (%v) source factory id, returned (%v)", action, ContainerSourceFactoryID, check)
		}
		if check := provider.params.decoderFactoryID; check != ContainerDecoderFactoryID {
			t.Errorf("%s didn't store the expected (%v) decoder factory id, returned (%v)", action, ContainerDecoderFactoryID, check)
		}
		if check := provider.params.observeFrequency; check != ContainerObserveFrequency {
			t.Errorf("%s didn't store the expected (%v) observe frequency, returned (%v)", action, ContainerObserveFrequency, check)
		}
		if check := provider.params.baseSourceID; check != ContainerBaseSourceID {
			t.Errorf("%s didn't store the expected (%v) base source id, returned (%v)", action, ContainerBaseSourceID, check)
		}
		if check := provider.params.baseSourcePath; check != ContainerBaseSourcePath {
			t.Errorf("%s didn't store the expected (%v) base source path, returned (%v)", action, ContainerBaseSourcePath, check)
		}
		if check := provider.params.baseSourceFormat; check != ContainerBaseSourceFormat {
			t.Errorf("%s didn't store the expected (%v) base source format, returned (%v)", action, ContainerBaseSourceFormat, check)
		}
	})

	t.Run("should creates a new config provider with given parameters", func(t *testing.T) {
		action := "Creating a new config provider with parameters"

		id := "__dummy_id__"
		fileSystemID := "__dummy_file_system_id__"
		loaderID := "__dummy_cloader_id__"
		sourceFactoryID := "__dummy_source_factory_id__"
		decoderFactoryID := "__dummy_decoder_factory_id__"
		observeFrequency := time.Second * 11
		baseSourceID := "__dummy_base_source_id__"
		baseSourcePath := "__dummy_base_source_path__"
		baseSourceFormat := "__dummy_base_source_format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		parameters := NewMockProviderParameters(ctrl)
		parameters.EXPECT().GetID().Return(id).Times(1)
		parameters.EXPECT().GetFileSystemID().Return(fileSystemID).Times(1)
		parameters.EXPECT().GetLoaderID().Return(loaderID).Times(1)
		parameters.EXPECT().GetSourceFactoryID().Return(sourceFactoryID).Times(1)
		parameters.EXPECT().GetDecoderFactoryID().Return(decoderFactoryID).Times(1)
		parameters.EXPECT().GetObserveFrequency().Return(observeFrequency).Times(1)
		parameters.EXPECT().GetBaseSourceID().Return(baseSourceID).Times(1)
		parameters.EXPECT().GetBaseSourcePath().Return(baseSourcePath).Times(1)
		parameters.EXPECT().GetBaseSourceFormat().Return(baseSourceFormat).Times(1)

		provider := NewProvider(parameters).(*provider)
		if provider == nil {
			t.Errorf("%s didn't return a valid reference to a new config provider", action)
		}
		if check := provider.params.id; check != id {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, id, check)
		}
		if check := provider.params.fileSystemID; check != fileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, fileSystemID, check)
		}
		if check := provider.params.loaderID; check != loaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, loaderID, check)
		}
		if check := provider.params.sourceFactoryID; check != sourceFactoryID {
			t.Errorf("%s didn't store the expected (%v) source factory id, returned (%v)", action, sourceFactoryID, check)
		}
		if check := provider.params.decoderFactoryID; check != decoderFactoryID {
			t.Errorf("%s didn't store the expected (%v) decoder factory id, returned (%v)", action, decoderFactoryID, check)
		}
		if check := provider.params.observeFrequency; check != observeFrequency {
			t.Errorf("%s didn't store the expected (%v) observe frequency, returned (%v)", action, observeFrequency, check)
		}
		if check := provider.params.baseSourceID; check != baseSourceID {
			t.Errorf("%s didn't store the expected (%v) base source id, returned (%v)", action, baseSourceID, check)
		}
		if check := provider.params.baseSourcePath; check != baseSourcePath {
			t.Errorf("%s didn't store the expected (%v) base source path, returned (%v)", action, baseSourcePath, check)
		}
		if check := provider.params.baseSourceFormat; check != baseSourceFormat {
			t.Errorf("%s didn't store the expected (%v) base source path, returned (%v)", action, baseSourceFormat, check)
		}
	})
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should retrieve the decoder factory", func(t *testing.T) {
		action := "calling the decoder factory builded by the registed provider"

		decoderFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, decoderFactoryMatcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := decoderFactoryMatcher.factory(container)
		switch entity.(type) {
		case DecoderFactory:
		default:
			t.Errorf("%s didn't return a valid reference to a new decoder factory", action)
		}
	})

	t.Run("should panic when trying to retrieve a source factory when the file system adapter is missing", func(t *testing.T) {
		action := "calling the source factory builded by the registed provider when the file system adapter is missing"

		sourceFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, sourceFactoryMatcher),
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
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerSourceFactoryID)
			}
		}()

		sourceFactoryMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a source factory and file system adapter is non-compliant with the interface", func(t *testing.T) {
		action := "calling the source factory builded by the registed provider when the file system adapter is non-compliant with the interface"

		sourceFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, sourceFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerSourceFactoryID)
			}
		}()

		sourceFactoryMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a source factory when the decoder factory is missing", func(t *testing.T) {
		action := "calling the source factory builded by the registed provider when the decoder factory is missing"

		sourceFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, sourceFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerSourceFactoryID)
			}
		}()

		sourceFactoryMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a source factory when the decoder factory is non-compliant with the interface", func(t *testing.T) {
		action := "calling the source factory builded by the registed provider when the decoder factory is non-compliant with the interface"

		sourceFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, sourceFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerSourceFactoryID)
			}
		}()

		sourceFactoryMatcher.factory(container)
	})

	t.Run("should retrieve the source factory", func(t *testing.T) {
		action := "calling the source factory builded by the registed provider"

		sourceFactoryMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, sourceFactoryMatcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(sys.ContainerFileSystemID).Return(fileSystem),
			container.EXPECT().Get(ContainerDecoderFactoryID).Return(decoderFactory),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := sourceFactoryMatcher.factory(container)
		switch entity.(type) {
		case SourceFactory:
		default:
			t.Errorf("%s didn't return a valid reference to a new source factory", action)
		}
	})

	t.Run("should retrieve the config", func(t *testing.T) {
		action := "calling the config builded by the registed provider"

		configMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, configMatcher),
			container.EXPECT().Add(ContainerLoaderID, matcher),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := configMatcher.factory(container)
		switch entity.(type) {
		case Config:
		default:
			t.Errorf("%s didn't return a valid reference to a new config", action)
		}
	})

	t.Run("should panic when trying to retrieve a loader when the config is missing", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the config is missing"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the config is non-compliant with the interface", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the config is non-compliant with the interface"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the source factory is missing", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the source factory is missing"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return(nil),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerID)
			}
		}()

		loaderMatcher.factory(container)
	})

	t.Run("should panic when trying to retrieve a loader when the source factory is non-compliant with the interface", func(t *testing.T) {
		action := "calling the loader builded by the registed provider when the source factory is non-compliant with the interface"

		loaderMatcher := &containerFactoryMatcher{}
		matcher := &containerFactoryMatcher{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return("__something_else__"),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s when requesting (%s) did not panic", action, ContainerID)
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

		config := NewMockConfig(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)

		container := NewMockContainer(ctrl)
		gomock.InOrder(
			container.EXPECT().Add(ContainerDecoderFactoryID, matcher),
			container.EXPECT().Add(ContainerSourceFactoryID, matcher),
			container.EXPECT().Add(ContainerID, matcher),
			container.EXPECT().Add(ContainerLoaderID, loaderMatcher),
		)
		gomock.InOrder(
			container.EXPECT().Get(ContainerID).Return(config),
			container.EXPECT().Get(ContainerSourceFactoryID).Return(sourceFactory),
		)

		provider := NewProvider(nil).(*provider)
		provider.Register(container)

		entity := loaderMatcher.factory(container)
		switch entity.(type) {
		case Loader:
		default:
			t.Errorf("%s didn't return a valid reference to a new loader", action)
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

	t.Run("should load the configuration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		loader := NewMockLoader(ctrl)
		loader.EXPECT().Load(ContainerBaseSourceID, ContainerBaseSourcePath, ContainerBaseSourceFormat).Times(1)

		container := NewMockContainer(ctrl)
		container.EXPECT().Get(ContainerLoaderID).Return(loader)

		provider := NewProvider(nil).(*provider)
		provider.Boot(container)
	})
}
