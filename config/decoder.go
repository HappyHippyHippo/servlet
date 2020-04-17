package config

// Decoder interface defines the intraction methods to a config content decoder
// used to parse the source content into a application usable configuration
// partial instance.
type Decoder interface {
	Close() error
	Decode() (Partial, error)
}
