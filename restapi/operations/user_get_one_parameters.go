// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewUserGetOneParams creates a new UserGetOneParams object
// no default values defined in spec.
func NewUserGetOneParams() UserGetOneParams {

	return UserGetOneParams{}
}

// UserGetOneParams contains all the bound params for the user get one operation
// typically these are obtained from a http.Request
//
// swagger:parameters userGetOne
type UserGetOneParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Идентификатор пользователя.
	  Required: true
	  In: path
	*/
	Nickname string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUserGetOneParams() beforehand.
func (o *UserGetOneParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rNickname, rhkNickname, _ := route.Params.GetOK("nickname")
	if err := o.bindNickname(rNickname, rhkNickname, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindNickname binds and validates parameter Nickname from path.
func (o *UserGetOneParams) bindNickname(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Nickname = raw

	return nil
}
