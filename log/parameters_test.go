package log

import (
	"os"
	"testing"

	"github.com/happyhippyhippo/servlet/config"
	"github.com/happyhippyhippo/servlet/sys"
)

func Test_NewParameters(t *testing.T) {
	loggerID := "logger_id"
	fileSystemID := "file_system_id"
	configID := "config_id"
	formatterFactoryID := "formatter_factory_id"
	streamFactoryID := "stream_factory_id"
	loaderID := "loader_id"

	t.Run("creates a new parameters", func(t *testing.T) {
		parameters := NewParameters(
			loggerID,
			fileSystemID,
			configID,
			formatterFactoryID,
			streamFactoryID,
			loaderID,
		)

		if value := parameters.LoggerID; value != loggerID {
			t.Errorf("stored (%v) logger ID", value)
		} else if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FormatterFactoryID; value != formatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		} else if value := parameters.StreamFactoryID; value != streamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		} else if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})
}

func Test_NewDefaultParameters(t *testing.T) {
	t.Run("creates a new parameters with the default values", func(t *testing.T) {
		parameters := NewDefaultParameters()

		if value := parameters.LoggerID; value != ContainerLoggerID {
			t.Errorf("stored (%v) logger ID", value)
		} else if value := parameters.FileSystemID; value != sys.ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.ConfigID; value != config.ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FormatterFactoryID; value != ContainerFormatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		} else if value := parameters.StreamFactoryID; value != ContainerStreamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})

	loggerID := "logger_id"
	fileSystemID := "file_system_id"
	configID := "config_id"
	formatterFactoryID := "formatter_factory_id"
	streamFactoryID := "stream_factory_id"
	loaderID := "loader_id"

	t.Run("creates a new parameters with the env logger ID", func(t *testing.T) {
		os.Setenv(EnvContainerLoggerID, loggerID)
		defer os.Setenv(EnvContainerLoggerID, "")

		parameters := NewDefaultParameters()
		if value := parameters.LoggerID; value != loggerID {
			t.Errorf("stored (%v) logger ID", value)
		}
	})

	t.Run("creates a new parameters with the env file system ID", func(t *testing.T) {
		os.Setenv(sys.EnvContainerFileSystemID, fileSystemID)
		defer os.Setenv(sys.EnvContainerFileSystemID, "")

		parameters := NewDefaultParameters()
		if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})

	t.Run("creates a new parameters with the env config ID", func(t *testing.T) {
		os.Setenv(config.EnvContainerConfigID, configID)
		defer os.Setenv(config.EnvContainerConfigID, "")

		parameters := NewDefaultParameters()
		if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		}
	})

	t.Run("creates a new parameters with the env formatter factory ID", func(t *testing.T) {
		os.Setenv(EnvContainerFormatterFactoryID, formatterFactoryID)
		defer os.Setenv(EnvContainerFormatterFactoryID, "")

		parameters := NewDefaultParameters()
		if value := parameters.FormatterFactoryID; value != formatterFactoryID {
			t.Errorf("stored (%v) formatter factory ID", value)
		}
	})

	t.Run("creates a new parameters with the env stream factory ID", func(t *testing.T) {
		os.Setenv(EnvContainerStreamFactoryID, streamFactoryID)
		defer os.Setenv(EnvContainerStreamFactoryID, "")

		parameters := NewDefaultParameters()
		if value := parameters.StreamFactoryID; value != streamFactoryID {
			t.Errorf("stored (%v) stream factory ID", value)
		}
	})

	t.Run("creates a new parameters with the env loader ID", func(t *testing.T) {
		os.Setenv(EnvContainerLoaderID, loaderID)
		defer os.Setenv(EnvContainerLoaderID, "")

		parameters := NewDefaultParameters()
		if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})
}
