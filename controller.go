package servlet

// Controller interface defines the interaction methods of a
// application controller used to process an inbound HTTP connection request.
type Controller interface {
	Options(context Context)
	Head(context Context)
	Get(context Context)
	Post(context Context)
	Put(context Context)
	Patch(context Context)
	Delete(context Context)
}

type controller struct{}

// Options method would be called on a OPTIONS HTTP verb request.
func (c controller) Options(context Context) {
	context.String(405, "")
}

// Head method would be called on a HEAD HTTP verb request.
func (c controller) Head(context Context) {
	context.String(405, "")
}

// Get method would be called on a GET HTTP verb request.
func (c controller) Get(context Context) {
	context.String(405, "")
}

// Post method would be called on a POST HTTP verb request.
func (c controller) Post(context Context) {
	context.String(405, "")
}

// Put method would be called on a PUT HTTP verb request.
func (c controller) Put(context Context) {
	context.String(405, "")
}

// Patch method would be called on a PATCH HTTP verb request.
func (c controller) Patch(context Context) {
	context.String(405, "")
}

// Delete method would be called on a DELETE HTTP verb request.
func (c controller) Delete(context Context) {
	context.String(405, "")
}
