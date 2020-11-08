package servlet

// ConfigDecoder interface defines the interaction methods to a config content
// decoder used to parse the source content into a application usable
// configuration partial instance.
type ConfigDecoder interface {
	Close()
	Decode() (ConfigPartial, error)
}
