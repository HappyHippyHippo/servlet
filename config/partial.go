package config

import (
	"strings"
)

// Partial interface define the methods to interact with a configuration
// partial section that defines a sub-section of the configuration coming
// from a particular source.
type Partial interface {
	Has(path string) bool
	Get(path string, def ...interface{}) interface{}
	Int(path string, def ...int) int
	String(path string, def ...string) string
	Config(path string, def ...Partial) Partial
}

type partial map[interface{}]interface{}

// Has will check if a requested path exists in the config partial.
func (p partial) Has(path string) bool {
	it := p
	nodes := strings.Split(path, ".")
	for i, node := range nodes {
		if node == "" {
			continue
		}

		switch it[node].(type) {
		case partial:
			it = it[node].(partial)
		case nil:
			return false
		default:
			return i == len(nodes)-1
		}
	}

	return true
}

// Get will retrieve the value stored in the requested path.
// If the path does not exists, then the value nil will be returned. Or, if
// a default value was given as the optional extra argument, then it will
// be returned instead of the standard nil value.
func (p partial) Get(path string, def ...interface{}) interface{} {
	it := p
	nodes := strings.Split(path, ".")
	for i, node := range nodes {
		if node == "" {
			continue
		}

		if _, ok := it[node]; !ok {
			if len(def) > 0 {
				return def[0]
			}
			return nil
		}

		switch it[node].(type) {
		case partial:
			it = it[node].(partial)
		case nil:
			return nil
		default:
			if i != len(nodes)-1 {
				if len(def) > 0 {
					return def[0]
				}
				return nil
			}
			return it[node]
		}
	}

	return it
}

// Int will return the casting to int of the stored value in the
// requested path. If the value retrieved was not found or returned nil, then
// the default optional argument will be returned if given.
func (p partial) Int(path string, def ...int) int {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return value.(int)
}

// String will return the casting to string of the stored value in the
// requested path. If the value retrieved was not found or returned nil, then
// the default optional argument will be returned if given.
func (p partial) String(path string, def ...string) string {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(string)
}

// Config will return the casting to a config partial of the stored
// value in the requested path. If the value retrieved was not found or
// returned nil, then the default optional argument will be returned if given.
func (p partial) Config(path string, def ...Partial) Partial {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(partial)
}

func (p partial) merge(p2 partial) partial {
	for key, value := range p2 {
		switch value.(type) {
		case partial:
			switch p[key].(type) {
			case partial:
				p[key] = p[key].(partial).merge(value.(partial))
			default:
				p[key] = value
			}
		default:
			p[key] = value
		}
	}
	return p
}
