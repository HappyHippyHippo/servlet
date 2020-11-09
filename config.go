package servlet

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

type configRefSource struct {
	id       string
	priority int
	source   ConfigSource
}

type configRefSourceSortByPriority []configRefSource

func (sources configRefSourceSortByPriority) Len() int {
	return len(sources)
}

func (sources configRefSourceSortByPriority) Swap(i, j int) {
	sources[i], sources[j] = sources[j], sources[i]
}

func (sources configRefSourceSortByPriority) Less(i, j int) bool {
	return sources[i].priority < sources[j].priority
}

type configRefObserver struct {
	path     string
	current  interface{}
	callback ConfigObserver
}

// Config defines the instance of a configuration managing structure.
type Config struct {
	mutex     sync.Locker
	sources   []configRefSource
	observers []configRefObserver
	partial   ConfigPartial
	loader    *TriggerRecurring
}

// NewConfig instantiate a new configuration object.
// This object will manage a series of sources, along side of the ability of
// registration of configuration path/values observer callbacks that will be
// called whenever the value has changed.
func NewConfig(period time.Duration) (*Config, error) {
	var c *Config

	var loader *TriggerRecurring
	if period != 0 {
		loader, _ = NewTriggerRecurring(period, func() error { return c.reload() })
	}

	c = &Config{
		mutex:     &sync.Mutex{},
		sources:   []configRefSource{},
		observers: []configRefObserver{},
		partial:   ConfigPartial{},
		loader:    loader,
	}

	return c, nil
}

// Close terminates the config instance.
// This will stop the observer trigger and call close on all registered sources.
func (c *Config) Close() {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if c.loader != nil {
		c.loader.Stop()
		c.loader = nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		reg.source.Close()
	}
}

// Has will check if a path has been loaded.
// This means that if the values has been loaded by any registered source.
func (c *Config) Has(path string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.partial.Has(path)
}

// Get will retrieve a configuration value loaded from a source.
func (c *Config) Get(path string, def ...interface{}) interface{} {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.partial.Get(path, def...)
}

// GetBool will retrieve a configuration value loaded from a
// source as a boolean.
func (c *Config) GetBool(path string, def ...bool) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case bool:
		return value.(bool)
	case string:
		v, err := strconv.ParseBool(value.(string))
		if err != nil {
			panic(err)
		}
		return v
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into bool", value, path))
}

// GetInt will retrieve a configuration value loaded from a source as a int.
func (c *Config) GetInt(path string, def ...int) int {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case int:
		return value.(int)
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 0)
		if err != nil {
			panic(err)
		}
		return int(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into int", value, path))
}

// GetInt8 will retrieve a configuration value loaded from a source as a int8.
func (c *Config) GetInt8(path string, def ...int8) int8 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case int8:
		return value.(int8)
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 8)
		if err != nil {
			panic(err)
		}
		return int8(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into int8", value, path))
}

// GetInt16 will retrieve a configuration value loaded from a source as a int16.
func (c *Config) GetInt16(path string, def ...int16) int16 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case int16:
		return value.(int16)
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 16)
		if err != nil {
			panic(err)
		}
		return int16(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into int16", value, path))
}

// GetInt32 will retrieve a configuration value loaded from a source as a int32.
func (c *Config) GetInt32(path string, def ...int32) int32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case int32:
		return value.(int32)
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			panic(err)
		}
		return int32(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into int32", value, path))
}

// GetInt64 will retrieve a configuration value loaded from a source as a int64.
func (c *Config) GetInt64(path string, def ...int64) int64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case int64:
		return value.(int64)
	case string:
		v, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into int64", value, path))
}

// GetUInt will retrieve a configuration value loaded from a source as a uint.
func (c *Config) GetUInt(path string, def ...uint) uint {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case uint:
		return value.(uint)
	case string:
		v, err := strconv.ParseUint(value.(string), 10, 0)
		if err != nil {
			panic(err)
		}
		return uint(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into uint", value, path))
}

// GetUInt8 will retrieve a configuration value loaded from a source as a uint8.
func (c *Config) GetUInt8(path string, def ...uint8) uint8 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case uint8:
		return value.(uint8)
	case string:
		v, err := strconv.ParseUint(value.(string), 10, 8)
		if err != nil {
			panic(err)
		}
		return uint8(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into uint8", value, path))
}

// GetUInt16 will retrieve a configuration value loaded from a
// source as a uint16.
func (c *Config) GetUInt16(path string, def ...uint16) uint16 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case uint16:
		return value.(uint16)
	case string:
		v, err := strconv.ParseUint(value.(string), 10, 16)
		if err != nil {
			panic(err)
		}
		return uint16(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into uint16", value, path))
}

// GetUInt32 will retrieve a configuration value loaded from a
// source as a uint32.
func (c *Config) GetUInt32(path string, def ...uint32) uint32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case uint32:
		return value.(uint32)
	case string:
		v, err := strconv.ParseUint(value.(string), 10, 32)
		if err != nil {
			panic(err)
		}
		return uint32(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into uint32", value, path))
}

// GetUInt64 will retrieve a configuration value loaded from a
// source as a uint64.
func (c *Config) GetUInt64(path string, def ...uint64) uint64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case uint64:
		return value.(uint64)
	case string:
		v, err := strconv.ParseUint(value.(string), 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into uint64", value, path))
}

// GetFloat32 will retrieve a configuration value loaded from a
// source as a float32.
func (c *Config) GetFloat32(path string, def ...float32) float32 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case float32:
		return value.(float32)
	case string:
		v, err := strconv.ParseFloat(value.(string), 32)
		if err != nil {
			panic(err)
		}
		return float32(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into float32", value, path))
}

// GetFloat64 will retrieve a configuration value loaded from a
// source as a float64.
func (c *Config) GetFloat64(path string, def ...float64) float64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case float64:
		return value.(float64)
	case string:
		v, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			panic(err)
		}
		return v
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into float64", value, path))
}

// GetComplex64 will retrieve a configuration value loaded from a
// source as a complex64.
func (c *Config) GetComplex64(path string, def ...complex64) complex64 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case complex64:
		return value.(complex64)
	case string:
		v, err := strconv.ParseComplex(value.(string), 64)
		if err != nil {
			panic(err)
		}
		return complex64(v)
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into complex64", value, path))
}

// GetComplex128 will retrieve a configuration value loaded from a
// source as a complex128.
func (c *Config) GetComplex128(path string, def ...complex128) complex128 {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case complex128:
		return value.(complex128)
	case string:
		v, err := strconv.ParseComplex(value.(string), 128)
		if err != nil {
			panic(err)
		}
		return v
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into complex128", value, path))
}

// GetRune will retrieve a configuration value loaded from a source as a rune.
func (c *Config) GetRune(path string, def ...rune) rune {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	var value interface{}
	if len(def) > 0 {
		value = c.Get(path, def[0])
	} else {
		value = c.Get(path)
	}

	switch value.(type) {
	case nil:
		panic(fmt.Errorf("path (%s) not found", path))
	case rune:
		return value.(rune)
	case string:
		return []rune(value.(string))[0]
	default:
	}

	panic(fmt.Errorf("unable to convert (%v) from path (%s) into rune", value, path))
}

// GetString will retrieve a configuration value loaded from a
// source as a string.
func (c *Config) GetString(path string, def ...string) string {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if len(def) > 0 {
		return c.Get(path, def[0]).(string)
	}
	return c.Get(path).(string)
}

// HasSource check if a source with a specific id has been registered.
func (c *Config) HasSource(id string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return true
		}
	}
	return false
}

// AddSource register a new source with a specific id with a given priority.
func (c *Config) AddSource(id string, priority int, source ConfigSource) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if source == nil {
		return fmt.Errorf("invalid nil 'source' argument")
	}
	if c.HasSource(id) {
		return fmt.Errorf("duplicate source id : %s", id)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.sources = append(c.sources, configRefSource{id, priority, source})
	sort.Sort(configRefSourceSortByPriority(c.sources))
	c.rebuild()

	return nil
}

// RemoveSource remove a source from the registration list
// of the configuration. This will also update the configuration content and
// re-validate the observed paths.
func (c *Config) RemoveSource(id string) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

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

// Source retrieve a previously registered source with a requested id.
func (c *Config) Source(id string) (ConfigSource, error) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			return reg.source, nil
		}
	}
	return nil, fmt.Errorf("source not found : %s", id)
}

// SourcePriority set a priority value of a previously registered
// source with the specified id. This may change the defined values if there
// was a override process of the configuration paths of the changing source.
func (c *Config) SourcePriority(id string, priority int) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.sources {
		if reg.id == id {
			reg.priority = priority
			sort.Sort(configRefSourceSortByPriority(c.sources))
			c.rebuild()

			return nil
		}
	}
	return fmt.Errorf("source not found : %s", id)
}

// HasObserver check if there is a observer to a configuration value path.
func (c *Config) HasObserver(path string) bool {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, reg := range c.observers {
		if reg.path == path {
			return true
		}
	}
	return false
}

// AddObserver register a new observer to a configuration path.
func (c *Config) AddObserver(path string, callback ConfigObserver) error {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	if callback == nil {
		return fmt.Errorf("invalid nil 'callback' argument")
	}

	value := c.Get(path)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.observers = append(c.observers, configRefObserver{path, value, callback})

	return nil
}

// RemoveObserver remove a observer to a configuration path.
func (c *Config) RemoveObserver(path string) {
	if c == nil {
		panic(fmt.Errorf("nil pointer receiver"))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, reg := range c.observers {
		if reg.path == path {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			return
		}
	}
}

func (c *Config) reload() error {
	rebuild := false
	for _, ref := range c.sources {
		switch s := ref.source.(type) {
		case ConfigSourceObservable:
			changed, _ := s.Reload()
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

func (c *Config) rebuild() {
	p := ConfigPartial{}
	for _, reg := range c.sources {
		p.merge(reg.source.Get("").(ConfigPartial))
	}

	c.partial = p

	for _, observer := range c.observers {
		updated := c.partial.Get(observer.path)
		if !reflect.DeepEqual(observer.current, updated) {
			old := observer.current
			observer.current = updated

			observer.callback(old, updated)
		}
	}
}
