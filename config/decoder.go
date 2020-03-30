package config

// Decoder interface defines the intraction methods to a configuration decoder.
type Decoder interface {
	Close() error
	Decode() (Partial, error)
}
