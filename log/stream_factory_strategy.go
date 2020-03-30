package log

import "github.com/happyhippyhippo/servlet/config"

// StreamFactoryStrategy interface defines the methods of the stream
// factory strategy that can instantiate a particular stream.
type StreamFactoryStrategy interface {
	Accept(stype string, args ...interface{}) bool
	AcceptConfig(conf config.Partial) bool
	Create(args ...interface{}) (Stream, error)
	CreateConfig(conf config.Partial) (Stream, error)
}
