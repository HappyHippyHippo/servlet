package servlet

import (
	"strings"
)

// ConfigPartial defined a type used to store configuration information.
type ConfigPartial map[interface{}]interface{}

// Has will check if a requested path exists in the config partial.
func (p ConfigPartial) Has(path string) bool {
	it := p
	nodes := strings.Split(path, ".")
	for i, node := range nodes {
		if node == "" {
			continue
		}

		switch it[node].(type) {
		case ConfigPartial:
			it = it[node].(ConfigPartial)
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
func (p ConfigPartial) Get(path string, def ...interface{}) interface{} {
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
		case ConfigPartial:
			it = it[node].(ConfigPartial)
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
func (p ConfigPartial) Int(path string, def ...int) int {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return value.(int)
}

// String will return the casting to string of the stored value in the
// requested path. If the value retrieved was not found or returned nil, then
// the default optional argument will be returned if given.
func (p ConfigPartial) String(path string, def ...string) string {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(string)
}

// Config will return the casting to a config partial of the stored
// value in the requested path. If the value retrieved was not found or
// returned nil, then the default optional argument will be returned if given.
func (p ConfigPartial) Config(path string, def ...ConfigPartial) ConfigPartial {
	value := p.Get(path)
	if value == nil && len(def) > 0 {
		return def[0]
	}
	return p.Get(path).(ConfigPartial)
}

func (p ConfigPartial) merge(p2 ConfigPartial) ConfigPartial {
	for key, value := range p2 {
		switch value.(type) {
		case ConfigPartial:
			switch p[key].(type) {
			case ConfigPartial:
			default:
				p[key] = ConfigPartial{}
			}
			p[key].(ConfigPartial).merge(value.(ConfigPartial))
		default:
			p[key] = value
		}
	}
	return p
}
