// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ConstaConst/technopark-db-forum-api/models"
)

// PostUpdateOKCode is the HTTP code returned for type PostUpdateOK
const PostUpdateOKCode int = 200

/*PostUpdateOK Информация о сообщении.


swagger:response postUpdateOK
*/
type PostUpdateOK struct {

	/*
	  In: Body
	*/
	Payload *models.Post `json:"body,omitempty"`
}

// NewPostUpdateOK creates PostUpdateOK with default headers values
func NewPostUpdateOK() *PostUpdateOK {

	return &PostUpdateOK{}
}

// WithPayload adds the payload to the post update o k response
func (o *PostUpdateOK) WithPayload(payload *models.Post) *PostUpdateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post update o k response
func (o *PostUpdateOK) SetPayload(payload *models.Post) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostUpdateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostUpdateNotFoundCode is the HTTP code returned for type PostUpdateNotFound
const PostUpdateNotFoundCode int = 404

/*PostUpdateNotFound Сообщение отсутсвует в форуме.


swagger:response postUpdateNotFound
*/
type PostUpdateNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostUpdateNotFound creates PostUpdateNotFound with default headers values
func NewPostUpdateNotFound() *PostUpdateNotFound {

	return &PostUpdateNotFound{}
}

// WithPayload adds the payload to the post update not found response
func (o *PostUpdateNotFound) WithPayload(payload *models.Error) *PostUpdateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post update not found response
func (o *PostUpdateNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostUpdateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
