package servlet

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_NewConfigProviderParams(t *testing.T) {
	t.Run("new parameters", func(t *testing.T) {
		parameters := NewConfigProviderParams()

		if value := parameters.ConfigID; value != ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FileSystemID; value != ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.SourceFactoryStrategyFileID; value != ContainerConfigSourceFactoryStrategyFileID {
			t.Errorf("stored (%v) source factory strategy file ID", value)
		} else if value := parameters.SourceFactoryStrategyObservableFileID; value != ContainerConfigSourceFactoryStrategyObservableFileID {
			t.Errorf("stored (%v) source factory strategy observable file ID", value)
		} else if value := parameters.SourceFactoryStrategyEnvironmentID; value != ContainerConfigSourceFactoryStrategyEnvironmentID {
			t.Errorf("stored (%v) source factory strategy environment ID", value)
		} else if value := parameters.SourceFactoryID; value != ContainerConfigSourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		} else if value := parameters.DecoderFactoryStrategyYamlID; value != ContainerConfigDecoderFactoryStrategyYamlID {
			t.Errorf("stored (%v) decoder factory strategy yaml ID", value)
		} else if value := parameters.DecoderFactoryID; value != ContainerConfigDecoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerConfigLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		} else if value := parameters.ObserveFrequency; value != ConfigObserveFrequency {
			t.Errorf("stored (%v) observe frequecy", value)
		} else if value := parameters.EntrySourceActive; value != ConfigEntrySourceActive {
			t.Errorf("stored (%v) base source active", value)
		} else if value := parameters.EntrySourceID; value != ConfigEntrySourceID {
			t.Errorf("stored (%v) base source id", value)
		} else if value := parameters.EntrySourcePath; value != ConfigEntrySourcePath {
			t.Errorf("stored (%v) base source path", value)
		} else if value := parameters.EntrySourceFormat; value != ConfigEntrySourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})

	t.Run("with the env config ID", func(t *testing.T) {
		value := "config_id"
		_ = os.Setenv(EnvContainerConfigID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.ConfigID; check != value {
			t.Errorf("stored (%v) config ID", check)
		}
	})

	t.Run("with the env file system ID", func(t *testing.T) {
		value := "file_system_id"
		_ = os.Setenv(EnvContainerFileSystemID, value)
		defer func() { _ = os.Setenv(EnvContainerFileSystemID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.FileSystemID; check != value {
			t.Errorf("stored (%v) file system ID", check)
		}
	})

	t.Run("with the env source factory strategy file ID", func(t *testing.T) {
		value := "source_factory_strategy_id"
		_ = os.Setenv(EnvContainerConfigSourceFactoryStrategyFileID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigSourceFactoryStrategyFileID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.SourceFactoryStrategyFileID; check != value {
			t.Errorf("stored (%v) source factory strategy file ID", check)
		}
	})

	t.Run("with the env source factory strategy observable file ID", func(t *testing.T) {
		value := "source_factory_strategy_id"
		_ = os.Setenv(EnvContainerConfigSourceFactoryStrategyObservableFileID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigSourceFactoryStrategyObservableFileID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.SourceFactoryStrategyObservableFileID; check != value {
			t.Errorf("stored (%v) source factory strategy observable file ID", check)
		}
	})

	t.Run("with the env source factory strategy environment ID", func(t *testing.T) {
		value := "source_factory_strategy_id"
		_ = os.Setenv(EnvContainerConfigSourceFactoryStrategyEnvironmentID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigSourceFactoryStrategyEnvironmentID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.SourceFactoryStrategyEnvironmentID; check != value {
			t.Errorf("stored (%v) source factory strategy environment ID", check)
		}
	})

	t.Run("with the env source factory ID", func(t *testing.T) {
		value := "source_factory_id"
		_ = os.Setenv(EnvContainerConfigSourceFactoryID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigSourceFactoryID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.SourceFactoryID; check != value {
			t.Errorf("stored (%v) source factory ID", check)
		}
	})

	t.Run("with the env decoder factory strategy yaml ID", func(t *testing.T) {
		value := "decoder_factory_id"
		_ = os.Setenv(EnvContainerConfigDecoderFactoryStrategyYamlID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigDecoderFactoryStrategyYamlID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.DecoderFactoryStrategyYamlID; check != value {
			t.Errorf("stored (%v) decoder factory strategy yaml ID", check)
		}
	})

	t.Run("with the env decoder factory ID", func(t *testing.T) {
		value := "decoder_factory_id"
		_ = os.Setenv(EnvContainerConfigDecoderFactoryID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigDecoderFactoryID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.DecoderFactoryID; check != value {
			t.Errorf("stored (%v) decoder factory ID", check)
		}
	})

	t.Run("with the env loader ID", func(t *testing.T) {
		value := "loader_id"
		_ = os.Setenv(EnvContainerConfigLoaderID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigLoaderID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.LoaderID; check != value {
			t.Errorf("stored (%v) loader ID", check)
		}
	})

	t.Run("with the env observer frequency", func(t *testing.T) {
		value := time.Second * 10
		_ = os.Setenv(EnvConfigObserveFrequency, strconv.Itoa(int(value.Seconds())))
		defer func() { _ = os.Setenv(EnvConfigObserveFrequency, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.ObserveFrequency; check != value {
			t.Errorf("stored (%v) observe frequency", check)
		}
	})

	t.Run("with the env base source active", func(t *testing.T) {
		_ = os.Setenv(EnvConfigEntrySourceActive, fmt.Sprintf("%v", true))
		defer func() { _ = os.Setenv(EnvConfigEntrySourceActive, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.EntrySourceActive; !check {
			t.Errorf("stored (%v) base source ID", check)
		}
	})

	t.Run("with the env base source ID", func(t *testing.T) {
		value := "base_config_id"
		_ = os.Setenv(EnvConfigEntrySourceID, value)
		defer func() { _ = os.Setenv(EnvConfigEntrySourceID, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.EntrySourceID; check != value {
			t.Errorf("stored (%v) base source ID", check)
		}
	})

	t.Run("with the env base source path", func(t *testing.T) {
		value := "base_config_path"
		_ = os.Setenv(EnvConfigEntrySourcePath, value)
		defer func() { _ = os.Setenv(EnvConfigEntrySourcePath, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.EntrySourcePath; check != value {
			t.Errorf("stored (%v) base source path", check)
		}
	})

	t.Run("with the env base source format", func(t *testing.T) {
		value := "base_config_format"
		_ = os.Setenv(EnvConfigEntrySourceFormat, value)
		defer func() { _ = os.Setenv(EnvConfigEntrySourceFormat, "") }()

		parameters := NewConfigProviderParams()
		if check := parameters.EntrySourceFormat; check != value {
			t.Errorf("stored (%v) base source format", check)
		}
	})
}
