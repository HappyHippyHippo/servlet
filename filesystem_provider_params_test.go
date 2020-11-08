package servlet

import (
	"os"
	"testing"
)

func Test_NewFileSystemParams(t *testing.T) {
	t.Run("no env override", func(t *testing.T) {
		p := NewFileSystemProviderParams()
		if p.FileSystemID != ContainerFileSystemID {
			t.Errorf("stored the '%s' file system container id", p.FileSystemID)
		}
	})

	t.Run("with env override", func(t *testing.T) {
		value := "test_id"
		_ = os.Setenv(EnvContainerFileSystemID, value)
		defer func() { _ = os.Setenv(EnvContainerFileSystemID, "") }()

		p := NewFileSystemProviderParams()
		if check := p.FileSystemID; check != value {
			t.Errorf("stored the '%s' file system container id", check)
		}
	})
}
