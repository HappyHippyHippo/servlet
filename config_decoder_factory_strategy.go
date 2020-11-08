package servlet

// ConfigDecoderFactoryStrategy interface defines the methods of the decoder
// factory strategy that can validate creation requests and instantiation of a
// particular decoder.
type ConfigDecoderFactoryStrategy interface {
	Accept(format string, args ...interface{}) bool
	Create(args ...interface{}) (ConfigDecoder, error)
}
