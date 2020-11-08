package servlet

import (
	"os"
	"testing"
)

func Test_NewLogProviderParams(t *testing.T) {
	t.Run("new parameters", func(t *testing.T) {
		parameters := NewLogProviderParams()

		if value := parameters.LoggerID; value != ContainerLoggerID {
			t.Errorf("stored (%v) logger ID", value)
		} else if value := parameters.FileSystemID; value != ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.ConfigID; value != ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FormatterFactoryStrategyJSONID; value != ContainerLogFormatterFactoryStrategyJSONID {
			t.Errorf("stored (%v) formatter factory strategy json ID", value)
		} else if value := parameters.FormatterFactoryID; value != ContainerLogFormatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		} else if value := parameters.StreamFactoryStrategyFileID; value != ContainerLogStreamFactoryStrategyFileID {
			t.Errorf("stored (%v) stream factory strategy file ID", value)
		} else if value := parameters.StreamFactoryID; value != ContainerLogStreamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerLogLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})

	t.Run("with the env logger ID", func(t *testing.T) {
		value := "logger_id"
		_ = os.Setenv(EnvContainerLoggerID, value)
		defer func() { _ = os.Setenv(EnvContainerLoggerID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.LoggerID; check != value {
			t.Errorf("stored (%v) logger ID", check)
		}
	})

	t.Run("with the env file system ID", func(t *testing.T) {
		value := "file_system_id"
		_ = os.Setenv(EnvContainerFileSystemID, value)
		defer func() { _ = os.Setenv(EnvContainerFileSystemID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.FileSystemID; check != value {
			t.Errorf("stored (%v) file system ID", check)
		}
	})

	t.Run("with the env config ID", func(t *testing.T) {
		value := "config_id"
		_ = os.Setenv(EnvContainerConfigID, value)
		defer func() { _ = os.Setenv(EnvContainerConfigID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.ConfigID; check != value {
			t.Errorf("stored (%v) config ID", check)
		}
	})

	t.Run("with the env formatter factory strategy json ID", func(t *testing.T) {
		value := "formatter_factory_id"
		_ = os.Setenv(EnvContainerLogFormatterFactoryStrategyJSONID, value)
		defer func() { _ = os.Setenv(EnvContainerLogFormatterFactoryStrategyJSONID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.FormatterFactoryStrategyJSONID; check != value {
			t.Errorf("stored (%v) formatter factory strategy json ID", check)
		}
	})

	t.Run("with the env formatter factory ID", func(t *testing.T) {
		value := "formatter_factory_id"
		_ = os.Setenv(EnvContainerLogFormatterFactoryID, value)
		defer func() { _ = os.Setenv(EnvContainerLogFormatterFactoryID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.FormatterFactoryID; check != value {
			t.Errorf("stored (%v) formatter factory ID", check)
		}
	})

	t.Run("with the env stream factory strategy file ID", func(t *testing.T) {
		value := "stream_factory_id"
		_ = os.Setenv(EnvContainerLogStreamFactoryStrategyFileID, value)
		defer func() { _ = os.Setenv(EnvContainerLogStreamFactoryStrategyFileID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.StreamFactoryStrategyFileID; check != value {
			t.Errorf("stored (%v) stream factory strategy file ID", check)
		}
	})

	t.Run("with the env stream factory ID", func(t *testing.T) {
		value := "stream_factory_id"
		_ = os.Setenv(EnvContainerLogStreamFactoryID, value)
		defer func() { _ = os.Setenv(EnvContainerLogStreamFactoryID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.StreamFactoryID; check != value {
			t.Errorf("stored (%v) stream factory ID", check)
		}
	})

	t.Run("with the env loader ID", func(t *testing.T) {
		value := "loader_id"
		_ = os.Setenv(EnvContainerLogLoaderID, value)
		defer func() { _ = os.Setenv(EnvContainerLogLoaderID, "") }()

		parameters := NewLogProviderParams()
		if check := parameters.LoaderID; check != value {
			t.Errorf("stored (%v) loader ID", check)
		}
	})
}
