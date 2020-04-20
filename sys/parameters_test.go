package sys

import (
	"os"
	"testing"
)

func Test_NewParameters(t *testing.T) {
	id := "id"
	parameters := NewParameters(id)

	t.Run("creates a new parameters", func(t *testing.T) {
		if value := parameters.FileSystemID; value != id {
			t.Errorf("stored (%v) file system ID", value)
		}
	})
}

func Test_NewDefaultParameters(t *testing.T) {
	t.Run("creates a new parameters with the default values", func(t *testing.T) {
		parameters := NewDefaultParameters()
		if value := parameters.FileSystemID; value != ContainerFileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})

	t.Run("creates a new parameters with the file system adapter id environment override", func(t *testing.T) {
		fileSystemID := "id"
		os.Setenv(EnvContainerFileSystemID, fileSystemID)
		defer os.Setenv(EnvContainerFileSystemID, "")

		parameters := NewDefaultParameters()
		if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})
}
