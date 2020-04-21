package servlet

// ControllerInterface interface defines the interaction methods of a
// application controller used to process an inbound HTTP connection request.
type ControllerInterface interface {
	Options(context Context)
	Head(context Context)
	Get(context Context)
	Post(context Context)
	Put(context Context)
	Patch(context Context)
	Delete(context Context)
}

// Controller type defines the base structure to be used to implement a
// HTTP request handler
type Controller struct{}

// Options method would be called on a OPTIONS HTTP verb request.
func (c Controller) Options(context Context) {
	context.String(405, "")
}

// Head method would be called on a HEAD HTTP verb request.
func (c Controller) Head(context Context) {
	context.String(405, "")
}

// Get method would be called on a GET HTTP verb request.
func (c Controller) Get(context Context) {
	context.String(405, "")
}

// Post method would be called on a POST HTTP verb request.
func (c Controller) Post(context Context) {
	context.String(405, "")
}

// Put method would be called on a PUT HTTP verb request.
func (c Controller) Put(context Context) {
	context.String(405, "")
}

// Patch method would be called on a PATCH HTTP verb request.
func (c Controller) Patch(context Context) {
	context.String(405, "")
}

// Delete method would be called on a DELETE HTTP verb request.
func (c Controller) Delete(context Context) {
	context.String(405, "")
}
