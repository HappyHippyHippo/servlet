package sys

// ContainerFileSystemID defines the default id used to register the
// application file system adapter instance in the application container.
const ContainerFileSystemID = "servlet.filesystem"

// FileSystemProviderParameters interface defines the method to interact
// with the file system provider parameters container object.
type FileSystemProviderParameters interface {
	GetID() string
	SetID(id string)
}

type fileSystemProviderParameters struct {
	id string
}

// NewFileSystemProviderParameters will instantiate a new file system provider
// parameters object with the requested values.
func NewFileSystemProviderParameters(id string) FileSystemProviderParameters {
	return &fileSystemProviderParameters{
		id: id,
	}
}

// NewDefaultFileSystemProviderParameters instantiate a new file system
// provider parameters object with the default values.
func NewDefaultFileSystemProviderParameters() FileSystemProviderParameters {
	return &fileSystemProviderParameters{
		id: ContainerFileSystemID,
	}
}

// GetID getter of the id parameter value.
func (p fileSystemProviderParameters) GetID() string {
	return p.id
}

// SetID setter of the id parameter value.
func (p *fileSystemProviderParameters) SetID(id string) {
	p.id = id
}
