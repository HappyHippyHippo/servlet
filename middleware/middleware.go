package middleware

import (
	"github.com/happyhippyhippo/servlet"
)

// Middleware interface defines the methods of a instance that will be
// used as a gin framework middleware
type Middleware interface {
	Execute(context servlet.Context)
}
