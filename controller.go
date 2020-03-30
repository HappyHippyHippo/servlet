package servlet

// Controller interface defines the defined interaction method of a
// application controller.
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

// Options method would be called on a OPTIONS verb request.
func (c controller) Options(context Context) {
	context.String(405, "")
}

// Head method would be called on a HEAD verb request.
func (c controller) Head(context Context) {
	context.String(405, "")
}

// Get method would be called on a GET verb request.
func (c controller) Get(context Context) {
	context.String(405, "")
}

// Post method would be called on a POST verb request.
func (c controller) Post(context Context) {
	context.String(405, "")
}

// Put method would be called on a PUT verb request.
func (c controller) Put(context Context) {
	context.String(405, "")
}

// Patch method would be called on a PATCH verb request.
func (c controller) Patch(context Context) {
	context.String(405, "")
}

// Delete method would be called on a DELETE verb request.
func (c controller) Delete(context Context) {
	context.String(405, "")
}
