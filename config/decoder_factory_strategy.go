package config

// DecoderFactoryStrategy interface defines the methods of the decoder factory
// strategy that can instantiate a particular decoder.
type DecoderFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (Decoder, error)
}
