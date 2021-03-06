// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ThreadUpdateHandlerFunc turns a function with the right signature into a thread update handler
type ThreadUpdateHandlerFunc func(ThreadUpdateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ThreadUpdateHandlerFunc) Handle(params ThreadUpdateParams) middleware.Responder {
	return fn(params)
}

// ThreadUpdateHandler interface for that can handle valid thread update params
type ThreadUpdateHandler interface {
	Handle(ThreadUpdateParams) middleware.Responder
}

// NewThreadUpdate creates a new http.Handler for the thread update operation
func NewThreadUpdate(ctx *middleware.Context, handler ThreadUpdateHandler) *ThreadUpdate {
	return &ThreadUpdate{Context: ctx, Handler: handler}
}

/*ThreadUpdate swagger:route POST /thread/{slug_or_id}/details threadUpdate

Обновление ветки

Обновление ветки обсуждения на форуме.


*/
type ThreadUpdate struct {
	Context *middleware.Context
	Handler ThreadUpdateHandler
}

func (o *ThreadUpdate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewThreadUpdateParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
