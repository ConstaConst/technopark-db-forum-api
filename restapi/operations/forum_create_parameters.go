// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	models "github.com/ConstaConst/technopark-db-forum-api/models"
)

// NewForumCreateParams creates a new ForumCreateParams object
// no default values defined in spec.
func NewForumCreateParams() ForumCreateParams {

	return ForumCreateParams{}
}

// ForumCreateParams contains all the bound params for the forum create operation
// typically these are obtained from a http.Request
//
// swagger:parameters forumCreate
type ForumCreateParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Данные форума.
	  Required: true
	  In: body
	*/
	Forum *models.Forum
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewForumCreateParams() beforehand.
func (o *ForumCreateParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Forum
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("forum", "body"))
			} else {
				res = append(res, errors.NewParseError("forum", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Forum = &body
			}
		}
	} else {
		res = append(res, errors.Required("forum", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
