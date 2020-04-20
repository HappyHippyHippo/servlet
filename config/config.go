package config

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/happyhippyhippo/servlet/sys"
	"github.com/happyhippyhippo/servlet/trigger"
)

// Config interface defines the interaction methods of a configuration instance
// that is responsable to manage a list of configuration sources.
type Config interface {
	Close() error
	Has(path string) bool
	Get(path string, def ...interface{}) interface{}
	GetBool(path string, def ...bool) bool
	GetInt(path string, def ...int) int
	GetInt8(path string, def ...int8) int8
	GetInt16(path string, def ...int16) int16
	GetInt32(path string, def ...int32) int32
	GetInt64(path string, def ...int64) int64
	GetUInt(path string, def ...uint) uint
	GetUInt8(path string, def ...uint8) uint8
	GetUInt16(path string, def ...uint16) uint16
	GetUInt32(path string, def ...uint32) uint32
	GetUInt64(path string, def ...uint64) uint64
	GetFloat32(path string, def ...float32) float32
	GetFloat64(path string, def ...float64) float64
	GetComplex64(path string, def ...complex64) complex64
	GetComplex128(path string, def ...complex128) complex128
	GetRune(path string, def ...rune) rune
	GetString(path string, def ...string) string
	HasSource(id string) bool
	AddSource(id string, priority int, source Source) error
	RemoveSource(id string)
	Source(id string) (Source, error)
	SourcePriority(id string, priority int) error
	HasObserver(path string) bool
	AddObserver(path string, callback Observer) error
	RemoveObserver(path string)
}

type refSource struct {
	id       string
	priority int
	source   Source
}

type refSourceSortByPriority []refSource

func (sources refSourceSortByPriority) Len() int {
	return len(sources)
}

func (sources refSourceSortByPriority) Swap(i, j int) {
	sources[i], sources[j] = sources[j], sources[i]
}

func (sources refSourceSortByPriority) Less(i, j int) bool {
	return sources[i].priority < sources[j].priority
}

type refOberserver struct {
	path     string
	current  interface{}
	callback Observer
}

type config struct {
	mutex     sys.RWMutex
	sources   []refSource
	observers []refOberserver
	partial   Partial
	reloader  trigger.Trigger
}

// NewConfig instatiate a new configuration object.
// This object will manage a series of sources, along side of the ability of
// registration of configuration path/values observer callbacks that will be
// called whenever the value has changed.
func NewConfig(period time.Duration) (Config, error) {
	var c *config

	var reloader trigger.Trigger = nil
	if period != 0 {
		reloader, _ = trigger.NewRecurringTrigger(period, func() error {
			return c.reload()
		})
	}

	c = &config{
		mutex:     &sync.RWMutex{},
		sources:   []refSource{},
		observers: []refOberserver{},
		partial:   partial{},
		reloader:  reloader,
	}

	return c, nil
}

// Close terminates the config instance.
// This will stop the observer trigger and call close on all registed sources.
func (c *config) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.reloader != nil {
		c.reloader.Stop()
		c.reloader = nil
	}

	for _, reg := range c.sources {
		reg.source.Close()
	}

	return nil
}

// Has will check if a path has been loaded.
// This means that if the values has been loaded by any registed source.
func (c *config) Has(path string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.partial.Has(path)
}

// Get will retrieve a configuration value loaded from a source.
func (c *config) Get(path string, def ...interface{}) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.partial.Get(path, def...)
}

// GetBool will retrieve a configuration value loaded from a
// source as a boolean.
func (c *config) GetBool(path string, def ...bool) bool {
	if len(def) > 0 {
		return c.Get(path, def[0]).(bool)
	}
	return c.Get(path).(bool)
}

// GetInt will retrieve a configuration value loaded from a source as a int.
func (c *config) GetInt(path string, def ...int) int {
	if len(def) > 0 {
		return c.Get(path, def[0]).(int)
	}
	return c.Get(path).(int)
}

// GetInt8 will retrieve a configuration value loaded from a source as a int8.
func (c *config) GetInt8(path string, def ...int8) int8 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(int8)
	}
	return c.Get(path).(int8)
}

// GetInt16 will retrieve a configuration value loaded from a source as a int16.
func (c *config) GetInt16(path string, def ...int16) int16 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(int16)
	}
	return c.Get(path).(int16)
}

// GetInt32 will retrieve a configuration value loaded from a source as a int32.
func (c *config) GetInt32(path string, def ...int32) int32 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(int32)
	}
	return c.Get(path).(int32)
}

// GetInt64 will retrieve a configuration value loaded from a source as a int64.
func (c *config) GetInt64(path string, def ...int64) int64 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(int64)
	}
	return c.Get(path).(int64)
}

// GetUInt will retrieve a configuration value loaded from a source as a uint.
func (c *config) GetUInt(path string, def ...uint) uint {
	if len(def) > 0 {
		return c.Get(path, def[0]).(uint)
	}
	return c.Get(path).(uint)
}

// GetUInt8 will retrieve a configuration value loaded from a source as a uint8.
func (c *config) GetUInt8(path string, def ...uint8) uint8 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(uint8)
	}
	return c.Get(path).(uint8)
}

// GetUInt16 will retrieve a configuration value loaded from a
// source as a uint16.
func (c *config) GetUInt16(path string, def ...uint16) uint16 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(uint16)
	}
	return c.Get(path).(uint16)
}

// GetUInt32 will retrieve a configuration value loaded from a
// source as a uint32.
func (c *config) GetUInt32(path string, def ...uint32) uint32 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(uint32)
	}
	return c.Get(path).(uint32)
}

// GetUInt64 will retrieve a configuration value loaded from a
// source as a uint64.
func (c *config) GetUInt64(path string, def ...uint64) uint64 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(uint64)
	}
	return c.Get(path).(uint64)
}

// GetFloat32 will retrieve a configuration value loaded from a
// source as a float32.
func (c *config) GetFloat32(path string, def ...float32) float32 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(float32)
	}
	return c.Get(path).(float32)
}

// GetFloat64 will retrieve a configuration value loaded from a
// source as a float64.
func (c *config) GetFloat64(path string, def ...float64) float64 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(float64)
	}
	return c.Get(path).(float64)
}

// GetComplex64 will retrieve a configuration value loaded from a
// source as a complex64.
func (c *config) GetComplex64(path string, def ...complex64) complex64 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(complex64)
	}
	return c.Get(path).(complex64)
}

// GetComplex128 will retrieve a configuration value loaded from a
// source as a complex128.
func (c *config) GetComplex128(path string, def ...complex128) complex128 {
	if len(def) > 0 {
		return c.Get(path, def[0]).(complex128)
	}
	return c.Get(path).(complex128)
}

// GetRune will retrieve a configuration value loaded from a source as a rune.
func (c *config) GetRune(path string, def ...rune) rune {
	if len(def) > 0 {
		return c.Get(path, def[0]).(rune)
	}
	return c.Get(path).(rune)
}

// GetString will retrieve a configuration value loaded from a
// source as a string.
func (c *config) GetString(path string, def ...string) string {
	if len(def) > 0 {
		return c.Get(path, def[0]).(string)
	}
	return c.Get(path).(string)
}

// HasSource check if a source with a specific id has been registed.
func (c *config) HasSource(id string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return true
		}
	}
	return false
}

// AddSource register a new source with a specific id with a given priority.
func (c *config) AddSource(id string, priority int, source Source) error {
	if source == nil {
		return fmt.Errorf("Invalid nil 'source' argument")
	}
	if c.HasSource(id) {
		return fmt.Errorf("Duplicate source id : %s", id)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.sources = append(c.sources, refSource{id, priority, source})
	sort.Sort(refSourceSortByPriority(c.sources))
	c.rebuild()

	return nil
}

// RemoveSource remove a source from the registration list
// of the configuration. This will also update the configuration content and
// revalidate the observed paths.
func (c *config) RemoveSource(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, reg := range c.sources {
		if reg.id == id {
			reg.source.Close()
			c.sources = append(c.sources[:i], c.sources[i+1:]...)
			c.rebuild()
			return
		}
	}
}

// Source retrieve a previously registed source with a requested id.
func (c *config) Source(id string) (Source, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return reg.source, nil
		}
	}
	return nil, fmt.Errorf("Source not found : %s", id)
}

// SourcePriority set a priority value of a previously registed
// source with the specified id. This may change the defined values if there
// was a override process of the configuration paths of the changing source.
func (c *config) SourcePriority(id string, priority int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			reg.priority = priority
			sort.Sort(refSourceSortByPriority(c.sources))
			c.rebuild()

			return nil
		}
	}
	return fmt.Errorf("Source not found : %s", id)
}

// HasObserver check if there is a observer to a configuration value path.
func (c *config) HasObserver(path string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, reg := range c.observers {
		if reg.path == path {
			return true
		}
	}
	return false
}

// AddObserver register a new observer to a configuration path.
func (c *config) AddObserver(path string, callback Observer) error {
	if callback == nil {
		return fmt.Errorf("Invalid nil 'callback' argument")
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.observers = append(c.observers, refOberserver{path, c.Get(path), callback})

	return nil
}

// RemoveObserver remove a observer to a configuration path.
func (c *config) RemoveObserver(path string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for i, reg := range c.observers {
		if reg.path == path {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			return
		}
	}
}

func (c *config) reload() error {
	rebuild := false
	for _, ref := range c.sources {
		switch ref.source.(type) {
		case ObservableSource:
			changed, _ := ref.source.(ObservableSource).Reload()
			rebuild = rebuild || changed
		}
	}

	if rebuild {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.rebuild()
	}

	return nil
}

func (c *config) rebuild() {
	p := partial{}
	for _, reg := range c.sources {
		p = p.merge(reg.source.Get("").(partial))
	}

	c.partial = p

	for _, observer := range c.observers {
		new := c.partial.Get(observer.path)
		if !reflect.DeepEqual(observer.current, new) {
			old := observer.current
			observer.current = new

			observer.callback(old, new)
		}
	}
}
