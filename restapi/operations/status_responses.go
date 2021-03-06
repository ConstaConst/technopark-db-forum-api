// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/ConstaConst/technopark-db-forum-api/models"
)

// StatusOKCode is the HTTP code returned for type StatusOK
const StatusOKCode int = 200

/*StatusOK Кол-во записей в базе данных, включая помеченные как "удалённые".


swagger:response statusOK
*/
type StatusOK struct {

	/*
	  In: Body
	*/
	Payload *models.Status `json:"body,omitempty"`
}

// NewStatusOK creates StatusOK with default headers values
func NewStatusOK() *StatusOK {

	return &StatusOK{}
}

// WithPayload adds the payload to the status o k response
func (o *StatusOK) WithPayload(payload *models.Status) *StatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the status o k response
func (o *StatusOK) SetPayload(payload *models.Status) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *StatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
