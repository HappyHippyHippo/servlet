package sys

import "sync"

// RWMutex exposes the interface to a Read-Write-Mutex.
type RWMutex interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
	RLocker() sync.Locker
}
