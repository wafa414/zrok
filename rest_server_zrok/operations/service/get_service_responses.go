// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti-test-kitchen/zrok/rest_model_zrok"
)

// GetServiceOKCode is the HTTP code returned for type GetServiceOK
const GetServiceOKCode int = 200

/*
GetServiceOK ok

swagger:response getServiceOK
*/
type GetServiceOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model_zrok.Service `json:"body,omitempty"`
}

// NewGetServiceOK creates GetServiceOK with default headers values
func NewGetServiceOK() *GetServiceOK {

	return &GetServiceOK{}
}

// WithPayload adds the payload to the get service o k response
func (o *GetServiceOK) WithPayload(payload *rest_model_zrok.Service) *GetServiceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get service o k response
func (o *GetServiceOK) SetPayload(payload *rest_model_zrok.Service) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetServiceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetServiceUnauthorizedCode is the HTTP code returned for type GetServiceUnauthorized
const GetServiceUnauthorizedCode int = 401

/*
GetServiceUnauthorized unauthorized

swagger:response getServiceUnauthorized
*/
type GetServiceUnauthorized struct {
}

// NewGetServiceUnauthorized creates GetServiceUnauthorized with default headers values
func NewGetServiceUnauthorized() *GetServiceUnauthorized {

	return &GetServiceUnauthorized{}
}

// WriteResponse to the client
func (o *GetServiceUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// GetServiceNotFoundCode is the HTTP code returned for type GetServiceNotFound
const GetServiceNotFoundCode int = 404

/*
GetServiceNotFound not found

swagger:response getServiceNotFound
*/
type GetServiceNotFound struct {
}

// NewGetServiceNotFound creates GetServiceNotFound with default headers values
func NewGetServiceNotFound() *GetServiceNotFound {

	return &GetServiceNotFound{}
}

// WriteResponse to the client
func (o *GetServiceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetServiceInternalServerErrorCode is the HTTP code returned for type GetServiceInternalServerError
const GetServiceInternalServerErrorCode int = 500

/*
GetServiceInternalServerError internal server error

swagger:response getServiceInternalServerError
*/
type GetServiceInternalServerError struct {
}

// NewGetServiceInternalServerError creates GetServiceInternalServerError with default headers values
func NewGetServiceInternalServerError() *GetServiceInternalServerError {

	return &GetServiceInternalServerError{}
}

// WriteResponse to the client
func (o *GetServiceInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
