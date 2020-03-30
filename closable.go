package servlet

// Closable is the interface used to signal the container that
// the element should must be closed on removal.
type Closable interface {
	Close() error
}
